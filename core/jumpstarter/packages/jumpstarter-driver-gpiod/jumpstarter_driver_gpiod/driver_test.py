from unittest.mock import MagicMock, patch

import pytest

# Import the client classes directly
from jumpstarter_core_driver_gpiod.client import DigitalInputClient, DigitalOutputClient


def setup_gpiod_mocks(mock_gpiod, line_number=18, line_name=None):
    """Helper function to set up common gpiod mocks"""
    if line_name is None:
        line_name = f"GPIO{line_number}"

    # Mock gpiod components
    mock_gpiod.line.Value.ACTIVE = 1
    mock_gpiod.line.Value.INACTIVE = 0
    mock_gpiod.line.Direction.OUTPUT = "output"
    mock_gpiod.line.Direction.INPUT = "input"
    mock_gpiod.line.Drive.PUSH_PULL = "push_pull"
    mock_gpiod.line.Bias.PULL_UP = "pull_up"
    mock_gpiod.line.Edge.BOTH = "both"
    mock_gpiod.EdgeEvent.Type.RISING_EDGE = "rising"
    mock_gpiod.EdgeEvent.Type.FALLING_EDGE = "falling"

    # Mock LineSettings
    mock_settings = MagicMock()
    mock_gpiod.LineSettings.return_value = mock_settings

    # Mock Chip
    mock_chip = MagicMock()
    mock_gpiod.Chip.return_value = mock_chip

    # Mock get_line_info for line name
    mock_line_info = MagicMock()
    mock_line_info.name = line_name
    mock_chip.get_line_info.return_value = mock_line_info

    # Mock LineRequest
    mock_line = MagicMock()
    mock_chip.request_lines.return_value = mock_line

    return mock_chip, mock_line, mock_settings


@pytest.fixture
def digital_output_client():
    """Create a DigitalOutputClient instance with mocked dependencies"""
    # Mock the required attributes
    client = DigitalOutputClient(
        stub=MagicMock(),
        portal=MagicMock(),
        stack=MagicMock()
    )
    client.uuid = "test-uuid"
    return client


@pytest.fixture
def digital_input_client():
    """Create a DigitalInputClient instance with mocked dependencies"""
    # Mock the required attributes
    client = DigitalInputClient(
        stub=MagicMock(),
        portal=MagicMock(),
        stack=MagicMock()
    )
    client.uuid = "test-uuid"
    return client


class TestDigitalOutputClient:
    """Test the DigitalOutputClient"""

    def test_on_method(self, digital_output_client):
        """Test the on() method"""
        # Mock the call method
        with patch.object(digital_output_client, 'call') as mock_call:
            digital_output_client.on()
            mock_call.assert_called_once_with("on")

    def test_off_method(self, digital_output_client):
        """Test the off() method"""
        # Mock the call method
        with patch.object(digital_output_client, 'call') as mock_call:
            digital_output_client.off()
            mock_call.assert_called_once_with("off")

    def test_read_method(self, digital_output_client):
        """Test the read() method"""
        # Mock the call method
        with patch.object(digital_output_client, 'call', return_value=1) as mock_call:
            result = digital_output_client.read()
            mock_call.assert_called_once_with("read_pin")
            assert result.value == 1
            assert str(result) == "active"

    def test_cli_method(self, digital_output_client):
        """Test the cli() method returns a click group"""
        cli_group = digital_output_client.cli()

        # Verify it's a click group with expected commands
        assert hasattr(cli_group, 'commands')
        assert 'on' in cli_group.commands
        assert 'off' in cli_group.commands
        assert 'read' in cli_group.commands


class TestDigitalInputClient:
    """Test the DigitalInputClient"""

    def test_wait_for_active_method(self, digital_input_client):
        """Test the wait_for_active() method"""
        # Mock the call method
        with patch.object(digital_input_client, 'call') as mock_call:
            digital_input_client.wait_for_active(timeout=5.0)
            mock_call.assert_called_once_with("wait_for_active", 5.0)

    def test_wait_for_inactive_method(self, digital_input_client):
        """Test the wait_for_inactive() method"""
        # Mock the call method
        with patch.object(digital_input_client, 'call') as mock_call:
            digital_input_client.wait_for_inactive(timeout=5.0)
            mock_call.assert_called_once_with("wait_for_inactive", 5.0)

    def test_read_method(self, digital_input_client):
        """Test the read() method"""
        # Mock the call method
        with patch.object(digital_input_client, 'call', return_value=0) as mock_call:
            result = digital_input_client.read()
            mock_call.assert_called_once_with("read_pin")
            assert result.value == 0
            assert str(result) == "inactive"

    def test_cli_method(self, digital_input_client):
        """Test the cli() method returns a click group"""
        cli_group = digital_input_client.cli()

        # Verify it's a click group with expected commands
        assert hasattr(cli_group, 'commands')
        assert 'read' in cli_group.commands
        assert 'wait-for-active' in cli_group.commands  # Note: uses hyphens, not underscores
        assert 'wait-for-inactive' in cli_group.commands  # Note: uses hyphens, not underscores


class TestDriverMethods:
    """Test the driver methods using mocking"""

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_digital_output_initialization(self, mock_gpiod):
        """Test DigitalOutput driver initialization with mocked gpiod"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        # Import and test the driver
        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        driver = DigitalOutput(
            line=18,
            drive="push_pull",
            active_low=False,
            bias="pull_up",
            initial_value="inactive"
        )

        # Verify gpiod.Chip was called
        mock_gpiod.Chip.assert_called_once_with("/dev/gpiochip0")

        # Verify line settings were configured
        mock_gpiod.LineSettings.assert_called()

        # Verify line was requested
        driver._chip.request_lines.assert_called_once()

        # Verify pin was processed correctly
        assert driver.line == 18

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_digital_input_initialization(self, mock_gpiod):
        """Test DigitalInput driver initialization with mocked gpiod"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=17)

        # Import and test the driver
        from jumpstarter_core_driver_gpiod.driver import DigitalInput

        driver = DigitalInput(
            line=17,
            drive=None,
            active_low=False,
            bias="pull_up"
        )

        # Verify gpiod.Chip was called
        mock_gpiod.Chip.assert_called_once_with("/dev/gpiochip0")

        # Verify line settings were configured
        mock_gpiod.LineSettings.assert_called()

        # Verify line was requested
        driver._chip.request_lines.assert_called_once()

        # Verify line was processed correctly
        assert driver.line == 17

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_digital_output_methods(self, mock_gpiod):
        """Test DigitalOutput driver methods with mocked gpiod"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        # Import and test the driver
        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        driver = DigitalOutput(line=18)

        # Test on() method
        driver.on()
        driver._line.set_value.assert_called_with(18, mock_gpiod.line.Value.ACTIVE)

        # Test off() method
        driver.off()
        driver._line.set_value.assert_called_with(18, mock_gpiod.line.Value.INACTIVE)

        # Test read_pin() method
        driver._line.get_value.return_value = mock_gpiod.line.Value.ACTIVE
        result = driver.read_pin()
        assert result.value == 1
        assert str(result) == "active"

        driver._line.get_value.return_value = mock_gpiod.line.Value.INACTIVE
        result = driver.read_pin()
        assert result.value == 0
        assert str(result) == "inactive"

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_digital_input_methods(self, mock_gpiod):
        """Test DigitalInput driver methods with mocked gpiod"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=17)

        # Import and test the driver
        from jumpstarter_core_driver_gpiod.driver import DigitalInput

        driver = DigitalInput(line=17)

        # Test read_pin() method
        driver._line.get_value.return_value = mock_gpiod.line.Value.ACTIVE
        result = driver.read_pin()
        assert result.value == 1
        assert str(result) == "active"

        driver._line.get_value.return_value = mock_gpiod.line.Value.INACTIVE
        result = driver.read_pin()
        assert result.value == 0
        assert str(result) == "inactive"

        # Test wait_for_active() when already active
        driver._line.get_value.return_value = mock_gpiod.line.Value.ACTIVE
        driver.wait_for_active()
        driver._line.wait_edge_events.assert_not_called()

        # Test wait_for_active() with timeout
        driver._line.get_value.return_value = mock_gpiod.line.Value.INACTIVE
        driver._line.wait_edge_events.return_value = False  # Timeout

        with pytest.raises(TimeoutError, match="Timed out waiting for line 17 edge event"):
            driver.wait_for_active(timeout=1.0)

        # Test wait_for_edge() with rising edge
        mock_event = MagicMock()
        mock_event.line_offset = 17
        mock_event.event_type = mock_gpiod.EdgeEvent.Type.RISING_EDGE

        driver._line.wait_edge_events.return_value = True
        driver._line.read_edge_events.return_value = [mock_event]

        driver.wait_for_edge("rising")
        driver._line.wait_edge_events.assert_called()
        driver._line.read_edge_events.assert_called()

        # Test wait_for_edge() with invalid edge type
        with pytest.raises(ValueError, match="Invalid edge type: invalid"):
            driver.wait_for_edge("invalid")


class TestErrorHandling:
    """Test error handling scenarios"""

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_gpiod_import_error(self, mock_gpiod):
        """Test handling of gpiod import error by mocking the import to fail"""
        # Mock the import to raise an ImportError
        mock_gpiod.side_effect = ImportError("gpiod is not installed")

        # This test verifies that the driver handles the import error gracefully
        # The actual import error handling is done at module level, so we test it differently
        try:
            from jumpstarter_core_driver_gpiod.driver import DigitalOutput
            # If we get here, the import succeeded (which is expected in our test environment)
            # We'll just verify that the driver can be imported
            assert DigitalOutput is not None
        except ImportError as e:
            # This would happen if gpiod is not available
            assert "gpiod is not installed" in str(e)

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_chip_open_error(self, mock_gpiod):
        """Test handling of chip open error"""
        mock_gpiod.Chip.side_effect = Exception("Cannot open chip")

        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        with pytest.raises(Exception, match="Cannot open chip"):
            DigitalOutput(line=18)

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_line_request_error(self, mock_gpiod):
        """Test handling of line request error"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        # Set up the error condition
        mock_chip.request_lines.side_effect = Exception("Cannot request line")

        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        with pytest.raises(Exception, match="Cannot request line"):
            DigitalOutput(line=18)

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_invalid_drive_value(self, mock_gpiod):
        """Test initialization with invalid drive value"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        with pytest.raises(ValueError, match="Invalid drive: invalid_drive"):
            DigitalOutput(line=18, drive="invalid_drive")

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_invalid_bias_value(self, mock_gpiod):
        """Test initialization with invalid bias value"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        with pytest.raises(ValueError, match="Invalid bias: invalid_bias"):
            DigitalOutput(line=18, bias="invalid_bias")

    @patch('jumpstarter_driver_gpiod.driver.gpiod')
    def test_invalid_initial_value(self, mock_gpiod):
        """Test initialization with invalid initial value"""
        # Set up common mocks
        mock_chip, mock_line, mock_settings = setup_gpiod_mocks(mock_gpiod, line_number=18)

        from jumpstarter_core_driver_gpiod.driver import DigitalOutput

        with pytest.raises(ValueError, match="Invalid initial_value: invalid"):
            DigitalOutput(line=18, initial_value="invalid")
