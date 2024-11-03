from django.db import models

class Order(models.Model):
    order_id = models.CharField(max_length=100, primary_key=True)
    operations = models.JSONField(default=list)

class OrderOperation(models.Model):
    guid = models.CharField(max_length=100)
    order = models.ForeignKey(Order, related_name='operations', on_delete=models.CASCADE)

class OrderOperationExecution(models.Model):
    guid = models.CharField(max_length=100)
    order = models.ForeignKey(Order, related_name='operation_executions', on_delete=models.CASCADE)

class OrderLabel(models.Model):
    label = models.CharField(max_length=100)
    order = models.ForeignKey(Order, on_delete=models.CASCADE)
    label_positions = models.JSONField(default=list)
    comment = models.TextField()