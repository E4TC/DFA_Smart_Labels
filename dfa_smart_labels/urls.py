"""dfa_smart_labels URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/3.2/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  path('', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  path('', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.urls import include, path
    2. Add a URL to urlpatterns:  path('blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import path
import drf_spectacular.views
import rest_framework.permissions
from rest_framework.routers import DefaultRouter
from django.urls import include


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
    
]
