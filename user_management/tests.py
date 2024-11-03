from django.test import TestCase


class EmptyTestDCase(TestCase):
    def test_empty(self) -> None:
        self.assertEqual(1, 1)

