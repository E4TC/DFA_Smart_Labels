from django.contrib.auth import forms as admin_forms
from django.forms import EmailField
from django.utils.translation import gettext_lazy as _

from user_management.models import User


class UserAdminChangeForm(admin_forms.UserChangeForm[User]):
    class Meta(admin_forms.UserChangeForm.Meta):  # type: ignore
        model = User
        field_classes = {"email": EmailField}


class UserAdminCreationForm(admin_forms.UserCreationForm[User]):
    class Meta(admin_forms.UserCreationForm.Meta):  # type: ignore
        model = User
        fields = ("email",)
        field_classes = {"email": EmailField}
        error_messages = {
            "email": {"unique": _("This email has already been taken.")},
        }
