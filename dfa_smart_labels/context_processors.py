import typing

import django.http
from django.conf import settings
from user_management.serializers import UserSerializer

DefaultContext = typing.TypedDict(
    "DefaultContext",
    {
        "static_url": str,
        "user_data": dict[str, typing.Any] | None,
        "version": str,
    },
)


def default(request: django.http.HttpRequest) -> DefaultContext:
    user_serialized = None
    if not request.user.is_anonymous:
        user_serialized = UserSerializer(request.user).data
    
    return {
        "static_url": settings.STATIC_URL,
        "user_data": user_serialized,
        "version": settings.APP_VERSION,
    }

