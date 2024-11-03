from django.contrib import admin
from django.urls import path
import drf_spectacular.views
import rest_framework.permissions
from rest_framework.routers import DefaultRouter
from django.urls import include

from backend.views import GetLabels, GetOrders, PostLabel


router = DefaultRouter()


class SpectacularAPIViewWithAuth(drf_spectacular.views.SpectacularAPIView):
    permission_classes = [rest_framework.permissions.IsAuthenticated]


class SpectacularSwaggerViewWithAuth(drf_spectacular.views.SpectacularSwaggerView):
    permission_classes = [rest_framework.permissions.IsAuthenticated]


urlpatterns = [
    path("admin/", admin.site.urls),
    path("api/", include(router.urls), name="api"),
    path("api/schema/", SpectacularAPIViewWithAuth.as_view(), name="schema"),
    path("api/docs/", SpectacularSwaggerViewWithAuth.as_view(url_name="schema"), name="swagger"),
    path('orders/', GetOrders.as_view(), name='get_orders'),
    path('label/', PostLabel.as_view(), name='post_label'),
    path('labels/', GetLabels.as_view(), name='get_labels'),
]
