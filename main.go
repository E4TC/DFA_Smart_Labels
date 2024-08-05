package main

import (
	"SmartLabels/models"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/influxdata/influxdb-client-go/v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// Store data from broker in variables, as we are not supposed to use a database
var order = models.Order{}
var orderOperation = models.OrderOperation{}
var orderOperationExecution = models.OrderOperationExecution{}

var orders []models.Order
var orderOperations []models.OrderOperation
var orderOperationExecutions []models.OrderOperationExecution

var hostname string

// Auth Token, valid for 10min, but we renew it every 5min
var sickToken models.SickToken

func main() {
	router := gin.Default()

	go checkToken()

	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	router.GET("/orders", getOrdersContext)
	router.POST("/label", postLabel)
	router.GET("/labels", getLabels)

	go router.Run(hostname)

	client := GetClient()
	if token := client.Subscribe("dfa/order/#", 0, onOrderReceived); token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error subscribing to topic:", token.Error())
	}
	client3 := GetClient("e4tc_operations_client")
	if token := client3.Subscribe("dfa/operation/#", 0, onOperationReceived); token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error subscribing to topic:", token.Error())
	}
	client2 := GetClient("e4tc_executions_client")
	if token := client2.Subscribe("dfa/operation_cycle/#", 0, onOperationExecutionReceived); token.Wait() && token.Error() != nil {
		fmt.Sprintf("Error subscribing to topic:", token.Error())
	}

	getDfaLabels()

	AnnounceService("e4tc_dfa_smartlabels")
	go StartHeartbeat("e4tc_dfa_smartlabels")
}

// We currently don't need the SICK API, but it might be useful in the future
func getSickToken() models.SickToken {
	var token models.SickToken

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "http://asset-analytics.io/api introspection")

	req, _ := http.NewRequest("POST", "https://192.168.205.226/user-manager/connect/token", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic ZTR0Yy1jbGllbnQ6ZTR0Y2NsaWVudA==")

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport, Timeout: 1 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Printf("impossible to send request: %s", err)
		return token
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return token
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return token
	}

	return token
}

// Get new SICK token every 5 minutes
func checkToken() {
	for {
		sickToken = getSickToken()
		time.Sleep(5 * time.Minute)
	}
}

// Return Orders list to API endpoint
func getOrdersContext(c *gin.Context) {
	var orders []models.Order
	orders = getOrders()
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// Get Data from API and send it to Bossard/Sepioo API and to the EDA broker
func postLabel(c *gin.Context) {
	if c.GetHeader("Token") == "z1s@2u9S86YN^KTpFS%^" {

		var newLabel models.OrderLabel
		if err := c.ShouldBindJSON(&newLabel); err != nil {
			return
		}

		statusCode := setLabelData(newLabel)
		newLabel.Timestamp = time.Now().Unix()
		newLabel.TimestampHr = time.Now().Format("2006-01-02 15:04:05")
		newLabel.Status = statusCode
		jsonObj, _ := json.Marshal(&newLabel)
		Publish("dfa/ot/labels/"+newLabel.Label, jsonObj)
		c.JSON(http.StatusOK, gin.H{"data": newLabel})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}

// Return Label list to API endpoint
func getLabels(c *gin.Context) {
	dfaLabels := getDfaLabels()
	c.JSON(http.StatusOK, gin.H{"data": dfaLabels})
}

// Get list of DFA Labels from Bossard/Sepioo API
func getDfaLabels() models.ActionUpdateLabelList {
	var allLabels models.ActionUpdateLabelList

	req, _ := http.NewRequest("GET", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("impossible to send request: %s", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	err = json.Unmarshal(body, &allLabels)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	// Filter DFA Labels
	var dfaLabels models.ActionUpdateLabelList
	for _, val := range allLabels {
		for _, tag := range val.Tags {
			if tag == "DFA" {
				dfaLabels = append(dfaLabels, val)
			}
		}
	}

	return dfaLabels
}

// When new order arrives from the broker add it to the array
func onOrderReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &order); err != nil {
		fmt.Println(err)
	}
	orders = append(orders, order)
}

// When new operation arrives from the broker add it to the array
func onOperationReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &orderOperation); err != nil {
		fmt.Println(err)
	}
	orderOperations = append(orderOperations, orderOperation)
}

// When new execution arrives from the broker add it to the array
func onOperationExecutionReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &orderOperationExecution); err != nil {
		fmt.Println(err)
	}
	orderOperationExecutions = append(orderOperationExecutions, orderOperationExecution)
}

// Use the latest data from the broker and create an Array of Orders which contains all associated operations and executions
func getOrders() []models.Order {
	orderMap := make(map[string]int)
	for i, order := range orders {
		orderMap[order.OrderID] = i
	}

	// Iterate through operations and assign them to matching orders using the map
	for _, operation := range orderOperations {
		if orderIndex, ok := orderMap[operation.Order]; ok {
			for i, oe := range orders[orderIndex].Operations {
				if oe.Guid == operation.Guid {
					orders[orderIndex].Operations = slices.Delete(orders[orderIndex].Operations, i, i+1)
				}
			}
			orders[orderIndex].Operations = append(orders[orderIndex].Operations, operation)
		}
	}

	// Iterate through operations and assign them to matching orders using the map
	for _, execution := range orderOperationExecutions {
		if orderIndex, ok := orderMap[execution.Order]; ok {
			for i, oe := range orders[orderIndex].OperationExecution {
				if oe.Guid == execution.Guid {
					orders[orderIndex].OperationExecution = slices.Delete(orders[orderIndex].OperationExecution, i, i+1)
				}
			}
			orders[orderIndex].OperationExecution = append(orders[orderIndex].OperationExecution, execution)
		}
	}
	return orders
}

// Use Bossard/Sepioo API to update Label Data
func setLabelData(data models.OrderLabel) int {
	update := models.UpdateLabel{Label: data.Label, Order: data.Order, LabelOperations: strings.Join(data.LabelOperations, " "), Comment: data.Comment}
	for id, value := range data.LabelPositions {
		switch id {
		case 0:
			update.Join1 = value
		case 1:
			update.Join2 = value
		case 2:
			update.Join3 = value
		case 3:
			update.Join4 = value
		case 4:
			update.Join5 = value
		case 5:
			update.Join6 = value
		case 6:
			update.Join7 = value
		}
	}
	allLabels := models.ActionUpdateLabel{Objects: data.Label, CustomFields: update}
	marshalled, _ := json.Marshal(allLabels)
	fmt.Println(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
	return res.StatusCode
}
