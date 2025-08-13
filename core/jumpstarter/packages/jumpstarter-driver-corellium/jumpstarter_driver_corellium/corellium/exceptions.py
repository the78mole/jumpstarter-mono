"""
Corellium API client exceptions module
"""
from jumpstarter_core.common.exceptions import JumpstarterException


class CorelliumApiException(JumpstarterException):
    """
    Exception raised when something goes wrong with Corellium's API.
    """
