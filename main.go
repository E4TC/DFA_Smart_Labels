package main

import (
	"SmartLabels/models"
	"bytes"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/influxdata/influxdb-client-go/v2"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"
)

var order = models.Order{}
var orderOperation = models.OrderOperation{}
var orderOperationExecution = models.OrderOperationExecution{}

var orders []models.Order
var orderOperations []models.OrderOperation
var orderOperationExecutions []models.OrderOperationExecution

func main() {
	router := gin.Default()

	models.ConnectDatabase()

	//go checkWorker()
	router.Use(static.Serve("/", static.LocalFile("./views", true)))

	router.GET("/locations", getLocations)
	router.POST("/locations", postLocations)
	router.GET("/locations/:id", findLocations)
	router.PATCH("/locations/:id", updateLocations)
	router.DELETE("/locations/:id", deleteLocations)

	router.GET("/workers", getWorkers)
	router.POST("/workers", postWorkers)
	router.GET("/workers/:id", findWorkers)
	router.PATCH("/workers/:id", updateWorkers)
	router.DELETE("/workers/:id", deleteWorkers)

	router.POST("/workers/:id/location/:lid/enter", enterWorkerLocation)
	router.POST("/workers/:id/location/:lid/leave", leaveWorkerLocation)

	router.POST("/button", pressButton)

	router.GET("/orders", getOrdersContext)
	router.POST("/label", postLabel)
	router.GET("/labels", getLabels)

	//router.Run("10.101.61.127:8080")
	go router.Run("localhost:8080")

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

	AnnounceService("e4tc_smartlabels")
	go StartHeartbeat("e4tc_smartlabels")

	//time.Sleep(time.Second * 3)
	//s, _ := json.MarshalIndent(getOrders(), "", "\t")
	//fmt.Println(string(s))

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	//client.Unsubscribe(topic)
	client.Disconnect(250)
}

func enterWorkerLocation(c *gin.Context) {
	var worker models.Workers
	if err := models.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
		return
	}

	var location models.Locations
	if err := models.DB.Where("id = ?", c.Param("lid")).First(&location).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location not found!"})
		return
	}

	// Display Data on the Labels
	//for _, pickOrder := range latestPickOrders {
	//	if !pickOrder.Done && pickOrder.Quantity != 0 {
	//		switchLabelPageAndQuantity(getLabelString(pickOrder.Text), pickOrder.Quantity)
	//	} else {
	//		switchLabelPage(getLabelString(pickOrder.Text), 3)
	//	}
	//}

	// Set Worker position to location
	var newLoc models.MoveWorkers
	newLoc.X = location.X
	newLoc.Y = location.Y
	newLoc.Z = location.Z
	models.DB.Model(&worker).Updates(newLoc)

}

func leaveWorkerLocation(c *gin.Context) {
	var worker models.Workers
	if err := models.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
		return
	}

	var location models.Locations
	if err := models.DB.Where("id = ?", c.Param("lid")).First(&location).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location not found!"})
		return
	}

	// Stop Flashing and Clear Label Screens
	stopAllLabelFlashing()
	switchAllLabelPage(1)

	var newLoc models.MoveWorkers
	newLoc.X = 1
	newLoc.Y = 1
	newLoc.Z = 1
	models.DB.Model(&worker).Updates(newLoc)
}

func checkWorker() {
	for {
		var workers []models.Workers
		if result := models.DB.Find(&workers).Error; result != nil {
			log.Fatal(result)
		}

		var locations []models.Locations
		if result := models.DB.Find(&locations).Error; result != nil {
			log.Fatal(result)
		}

		for _, worker := range workers {
			// Read worker location from external API
			response, err := http.Get(worker.URL)
			if err != nil {
				log.Print(err.Error())
				continue
			}
			responseData, err := io.ReadAll(response.Body)
			if err != nil {
				log.Fatal(err)
			}
			data := models.MoveWorkers{}
			json.Unmarshal([]byte(responseData), &data)

			// Calculate distance
			distance := float32(math.Sqrt(math.Pow(float64(worker.X-data.X), 2) + math.Pow(float64(worker.Y-data.Y), 2)))
			if distance > 0.3 { //ignore sensor noise
				data.Distance = worker.Distance + distance
			}

			// Update DB
			models.DB.Model(&worker).Updates(data)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func pressButton(c *gin.Context) {
	// Right now nothing happens on button press
}

func getLocations(c *gin.Context) {
	var locations []models.Locations
	models.DB.Find(&locations)
	c.JSON(http.StatusOK, gin.H{"data": locations})
}
func postLocations(c *gin.Context) {
	var newLocation models.Locations

	if err := c.ShouldBindJSON(&newLocation); err != nil {
		return
	}

	location := models.Locations{Text: newLocation.Text, X: newLocation.X, Y: newLocation.Y, Z: newLocation.Z}
	models.DB.Create(&location)

	c.JSON(http.StatusOK, gin.H{"data": location})
}
func findLocations(c *gin.Context) {
	var location models.Locations

	if err := models.DB.Where("id = ?", c.Param("id")).First(&location).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": location})
}
func updateLocations(c *gin.Context) {
	var location models.Locations

	if err := models.DB.Where("id = ?", c.Param("id")).First(&location).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location not found!"})
		return
	}

	var input models.UpdateLocations
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&location).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": location})
}
func deleteLocations(c *gin.Context) {
	var location models.Locations

	if err := models.DB.Where("id = ?", c.Param("id")).First(&location).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location not found!"})
		return
	}

	models.DB.Delete(&location)

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func getWorkers(c *gin.Context) {
	var workers []models.Workers
	models.DB.Find(&workers)
	c.JSON(http.StatusOK, gin.H{"data": workers})
}
func postWorkers(c *gin.Context) {
	if c.GetHeader("Token") == "z1s@2u9S86YN^KTpFS%^" {

		var newWorker models.Workers
		if err := c.ShouldBindJSON(&newWorker); err != nil {
			return
		}

		worker := models.Workers{Title: newWorker.Title, X: newWorker.X, Y: newWorker.Y, Z: newWorker.Z, Distance: newWorker.Distance}
		models.DB.Create(&worker)

		c.JSON(http.StatusOK, gin.H{"data": worker})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}
func findWorkers(c *gin.Context) {
	var worker models.Workers

	if err := models.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": worker})
}
func updateWorkers(c *gin.Context) {
	if c.GetHeader("Token") == "z1s@2u9S86YN^KTpFS%^" {
		var worker models.Workers

		if err := models.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
			return
		}

		var input models.UpdateWorkers
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		models.DB.Model(&worker).Updates(input)

		c.JSON(http.StatusOK, gin.H{"data": worker})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}
}
func deleteWorkers(c *gin.Context) {
	if c.GetHeader("Token") == "z1s@2u9S86YN^KTpFS%^" {
		var worker models.Workers

		if err := models.DB.Where("id = ?", c.Param("id")).First(&worker).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Worker not found!"})
			return
		}

		models.DB.Delete(&worker)

		c.JSON(http.StatusOK, gin.H{"data": true})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}
}

func getOrdersContext(c *gin.Context) {
	var orders []models.Order
	orders = getOrders()
	c.JSON(http.StatusOK, gin.H{"data": orders})
}
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
		Publish("dfa/labels/"+newLabel.Label, jsonObj)
		c.JSON(http.StatusOK, gin.H{"data": newLabel})
	} else {
		c.JSON(http.StatusForbidden, gin.H{})
	}

}
func getLabels(c *gin.Context) {
	dfaLabels := getDfaLabels()
	c.JSON(http.StatusOK, gin.H{"data": dfaLabels})
}

func stopAllLabelFlashing() {
	allLabels := models.ActionStopFlash{Objects: []string{"DemoObject1", "DemoObject2", "DemoObject3"}}
	marshalled, _ := json.Marshal(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects/action/stopflash", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
}
func switchLabelPage(label string, page uint) {
	allLabels := models.ActionSwitchPage{Objects: []string{label}, Page: page}
	marshalled, _ := json.Marshal(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects/action/switchpage", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
}
func switchAllLabelPage(page uint) {
	allLabels := models.ActionSwitchPage{Objects: []string{"DemoObject1", "DemoObject2", "DemoObject3"}, Page: page}
	marshalled, _ := json.Marshal(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects/action/switchpage", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
}
func flashLabelLED(label string) {
	allLabels := models.ActionFlash{Objects: []string{label}, Color: "RED", Duration: 5, Patter: "FLASH_1_SECOND"}
	marshalled, _ := json.Marshal(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects/action/flash", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
}
func pingAllLabels() {
	allLabels := models.ActionStopFlash{Objects: []string{"DemoObject1", "DemoObject2", "DemoObject3"}}
	marshalled, _ := json.Marshal(allLabels)
	req, _ := http.NewRequest("POST", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects/action/ping", bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
}

func getDfaLabels() models.ActionUpdateLabelList {
	var allLabels models.ActionUpdateLabelList

	req, _ := http.NewRequest("GET", "https://industrial-api.azure-api.net/v2.0/industry_e4tc_eu/E4TCDemo/objects", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Ocp-Apim-Subscription-Key", "3607c3cfaf414cb6bb8f24e57c10dd71")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	err = json.Unmarshal(body, &allLabels)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	//log.Printf("status Code: %d", res.StatusCode)

	// Filter DFA Labels
	var dfaLabels models.ActionUpdateLabelList
	for _, val := range allLabels {
		for _, tag := range val.Tags {
			if tag == "DFA" {
				dfaLabels = append(dfaLabels, val)
			}
		}
	}
	//fmt.Println(dfaLabels)
	return dfaLabels
}

func getLabelString(id uint) string {
	var label models.Locations
	if result := models.DB.Where("`id` = ?", id).Last(&label).Error; result != nil {
		log.Fatal(result)
	}
	return label.Text
}

func onOrderReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &order); err != nil {
		fmt.Println(err)
	}
	orders = append(orders, order)
}
func onOperationReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &orderOperation); err != nil {
		fmt.Println(err)
	}
	orderOperations = append(orderOperations, orderOperation)
}
func onOperationExecutionReceived(client mqtt.Client, message mqtt.Message) {
	if err := json.Unmarshal(message.Payload(), &orderOperationExecution); err != nil {
		fmt.Println(err)
	}
	orderOperationExecutions = append(orderOperationExecutions, orderOperationExecution)
}

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
func mapOperationToOrders(orders []models.Order, operation models.OrderOperation) {
	// Create a map to store orderID as key and index in the 'orders' slice as value
	orderMap := make(map[string]int)
	for i, order := range orders {
		orderMap[order.OrderID] = i
	}

	// assign operation to matching orders using the map
	if orderIndex, ok := orderMap[operation.Order]; ok {
		for i, oe := range orders[orderIndex].Operations {
			//fmt.Println(oe.Guid, operation.Guid)
			if oe.Guid == operation.Guid {
				orders[orderIndex].Operations = slices.Delete(orders[orderIndex].Operations, i, i+1)
			}
		}
		orders[orderIndex].Operations = append(orders[orderIndex].Operations, operation)
	}
}
func mapOperationExecutionToOrders(orders []models.Order, operationEx models.OrderOperationExecution) {
	// Create a map to store orderID as key and index in the 'orders' slice as value
	orderMap := make(map[string]int)
	for i, order := range orders {
		orderMap[order.OrderID] = i
	}

	// assign operation to matching orders using the map
	if orderIndex, ok := orderMap[operationEx.Order]; ok {
		for i, oe := range orders[orderIndex].OperationExecution {
			if oe.Guid == operationEx.Guid {
				orders[orderIndex].OperationExecution = slices.Delete(orders[orderIndex].OperationExecution, i, i+1)
			}
		}
		orders[orderIndex].OperationExecution = append(orders[orderIndex].OperationExecution, operationEx)
	}

}

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
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)
	return res.StatusCode
}
