from rest_framework import serializers
from models import (
    Order,
    OrderPosition,
    OrderOperation,
    OrderOperationExecution,
    OrderLabel,
    UpdateLabel,
    ActionUpdateLabelList,
    ActionUpdateLabel,
    OrderLabelPub
)


class OrderPositionSerializer(serializers.ModelSerializer):
    class Meta:
        model = OrderPosition
        fields = ['guid', 'position', 'description', 'quantity', 'state']


class OrderOperationSerializer(serializers.ModelSerializer):
    class Meta:
        model = OrderOperation
        fields = [
            'guid', 'order', 'position', 'operation_id', 'description', 
            'cost_group', 'machine_group', 'machine_name', 'state'
        ]


class OrderOperationExecutionSerializer(serializers.ModelSerializer):
    class Meta:
        model = OrderOperationExecution
        fields = [
            'guid', 'order', 'operation', 'operation_description', 
            'cost_group', 'machine_group', 'machine_name', 'start', 'stop'
        ]


class OrderLabelSerializer(serializers.ModelSerializer):
    class Meta:
        model = OrderLabel
        fields = [
            'label', 'order', 'label_operations', 'label_positions', 
            'comment', 'timestamp', 'timestamp_hr', 'status'
        ]


class UpdateLabelSerializer(serializers.ModelSerializer):
    class Meta:
        model = UpdateLabel
        fields = ['label', 'order', 'positions', 'operation_ids', 'comment']


class ActionUpdateLabelListSerializer(serializers.ModelSerializer):
    custom_fields = UpdateLabelSerializer()

    class Meta:
        model = ActionUpdateLabelList
        fields = ['object_id', 'custom_fields', 'tags']


class ActionUpdateLabelSerializer(serializers.ModelSerializer):
    custom_fields = UpdateLabelSerializer()

    class Meta:
        model = ActionUpdateLabel
        fields = ['object_id', 'custom_fields']


class OrderLabelPubSerializer(serializers.ModelSerializer):
    class Meta:
        model = OrderLabelPub
        fields = [
            'label', 'order', 'label_operations', 'label_positions', 
            'comment', 'timestamp', 'timestamp_hr'
        ]


class OrderSerializer(serializers.ModelSerializer):
    positions = OrderPositionSerializer(many=True, read_only=True)
    operations = OrderOperationSerializer(many=True, read_only=True)
    operation_executions = OrderOperationExecutionSerializer(many=True, read_only=True)

    class Meta:
        model = Order
        fields = [
            'order_id', 'description', 'customer', 'created', 
            'created_hr', 'positions', 'operations', 'operation_executions'
        ]
