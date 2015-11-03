from django.contrib import admin
from kombu.transport.django import models as kombu_models

# Register your models here.

admin.site.register(kombu_models.Message)
