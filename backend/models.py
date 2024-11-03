from django.db import models


class Order(models.Model):
    order_id = models.CharField(max_length=100, unique=True)
    description = models.TextField()
    customer = models.CharField(max_length=100)
    created = models.BigIntegerField()
    created_hr = models.CharField(max_length=100)

    def __str__(self):
        return self.order_id


class OrderPosition(models.Model):
    guid = models.CharField(max_length=100, unique=True)
    position = models.IntegerField()
    description = models.CharField(max_length=255)
    quantity = models.IntegerField()
    state = models.CharField(max_length=50)
    order = models.ForeignKey(Order, related_name='positions', on_delete=models.CASCADE)

    def __str__(self):
        return self.guid


class OrderOperation(models.Model):
    guid = models.CharField(max_length=100, unique=True)
    order = models.ForeignKey(Order, related_name='operations', on_delete=models.CASCADE)
    position = models.IntegerField()
    operation_id = models.IntegerField()
    description = models.TextField()
    cost_group = models.CharField(max_length=100)
    machine_group = models.CharField(max_length=100)
    machine_name = models.CharField(max_length=100)
    state = models.CharField(max_length=50)

    def __str__(self):
        return self.guid


class OrderOperationExecution(models.Model):
    guid = models.CharField(max_length=100, unique=True)
    order = models.ForeignKey(Order, related_name='operation_executions', on_delete=models.CASCADE)
    operation = models.CharField(max_length=100)
    operation_description = models.TextField()
    cost_group = models.CharField(max_length=100)
    machine_group = models.CharField(max_length=100)
    machine_name = models.CharField(max_length=100)
    start = models.CharField(max_length=100)
    stop = models.CharField(max_length=100)

    def __str__(self):
        return self.guid


class OrderLabel(models.Model):
    label = models.CharField(max_length=100)
    order = models.ForeignKey(Order, related_name='labels', on_delete=models.CASCADE)
    label_operations = models.JSONField()
    label_positions = models.JSONField() 
    comment = models.TextField()
    timestamp = models.BigIntegerField(null=True, blank=True)
    timestamp_hr = models.CharField(max_length=100, null=True, blank=True)
    status = models.IntegerField(null=True, blank=True)

    def __str__(self):
        return self.label


class UpdateLabel(models.Model):
    label = models.CharField(max_length=100)
    order = models.CharField(max_length=100)
    positions = models.CharField(max_length=100)
    operation_ids = models.JSONField()
    comment = models.TextField()

    def __str__(self):
        return self.label


class ActionUpdateLabelList(models.Model):
    object_id = models.CharField(max_length=100)
    custom_fields = models.ForeignKey(UpdateLabel, on_delete=models.CASCADE)
    tags = models.JSONField()

    def __str__(self):
        return self.object_id


class ActionUpdateLabel(models.Model):
    object_id = models.CharField(max_length=100)
    custom_fields = models.ForeignKey(UpdateLabel, on_delete=models.CASCADE)

    def __str__(self):
        return self.object_id


class OrderLabelPub(models.Model):
    label = models.CharField(max_length=100)
    order = models.CharField(max_length=100)
    label_operations = models.JSONField() 
    label_positions = models.JSONField()  
    comment = models.TextField()
    timestamp = models.BigIntegerField(null=True, blank=True)
    timestamp_hr = models.CharField(max_length=100, null=True, blank=True)

    def __str__(self):
        return self.label
