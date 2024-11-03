from typing import Any

from django.core.management.base import CommandParser
from django.core.management.commands.makemessages import Command as MMCommand


# Source: https://medium.com/@hugosousa/hacking-djangos-makemessages-for-better-translations-matching-in-jsx-components-1174b57a13b1 # noqa
class Command(MMCommand):
    def add_arguments(self, parser: CommandParser) -> None:
        parser.add_argument(
            "--language",
            "-lang",
            default="Python",
            dest="language",
            help="Language to be used by xgettext",
        )

        super(Command, self).add_arguments(parser)

    def handle(self, *args: list[Any], **options: Any) -> None:
        language = options.get("language")
        self.xgettext_options.append("--language={lang}".format(lang=language))
        super(Command, self).handle(*args, **options)
