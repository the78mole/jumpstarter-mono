from jumpstarter_core.common import exceptions


class LeaseError(exceptions.JumpstarterException):
    """Raised when a lease operation fails."""
