import json
import requests
from django.http import JsonResponse
from django.views import View
from .models import Order, OrderLabel

class GetOrders(View):
    def get(self, request):
        orders = list(Order.objects.all().values())
        return JsonResponse({"data": orders}, status=200)

class PostLabel(View):
    def post(self, request):
        try:
            new_label_data = json.loads(request.body)
            new_label = OrderLabel.objects.create(
                label=new_label_data['label'],
                order_id=new_label_data['order'],
                label_positions=new_label_data.get('label_positions', []),
                comment=new_label_data.get('comment', '')
            )
            return JsonResponse({"data": new_label_data}, status=201)
        except Exception as e:
            return JsonResponse({"error": str(e)}, status=400)

class GetLabels(View):
    def get(self, request):
        labels = list(OrderLabel.objects.all().values())
        return JsonResponse({"data": labels}, status=200)