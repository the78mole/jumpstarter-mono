# Hardware Documentation

This section covers the hardware components of the Jumpstarter system.

## Hardware Architecture Overview

```mermaid
graph TB
    subgraph "Host System"
        HOST[Host Computer<br/>Running Jumpstarter]
        USB[USB 3.0 Interface]
        ETH[Ethernet Interface]
    end

    subgraph "DUT Link Board"
        MCU[STM32 MCU<br/>Rust Firmware]
        POWER[Power Management<br/>3.3V, 5V, 12V]
        GPIO[16x GPIO Pins]
        UART[4x UART Ports]
        SPI[2x SPI Ports]
        I2C[2x I2C Ports]
        USB_PORTS[4x USB Ports]
        MONITOR[Current/Voltage<br/>Monitoring]
    end

    subgraph "DUT Pi-Link Board"
        PI[Raspberry Pi<br/>Linux System]
        RELAY[8-Channel Relay<br/>Waveshare Board]
        SDMUX[USB-SD-Mux FAST<br/>SD Card Flashing]
        SERIAL[USB-Serial<br/>Console Converter]
        ETH_USB[Ethernet-USB<br/>Dongle]
        PI_ETH[Pi Ethernet<br/>External Network]
        PI_GPIO[Pi GPIO<br/>Custom Control]
    end

    subgraph "Device Under Test"
        DUT[Target Device]
        POWER_IN[Power Input]
        DATA_IO[Data I/O]
        CONSOLE[Debug Console]
        SD_SLOT[SD Card Slot]
        DIP_SW[DIP Switches]
    end

    HOST --> USB
    HOST --> ETH
    USB --> MCU
    ETH --> MCU

    MCU --> POWER
    MCU --> GPIO
    MCU --> UART
    MCU --> SPI
    MCU --> I2C
    MCU --> USB_PORTS
    MCU --> MONITOR

    HOST --> PI
    PI --> RELAY
    PI --> SDMUX
    PI --> SERIAL
    PI --> ETH_USB
    PI --> PI_ETH
    PI --> PI_GPIO

    POWER --> POWER_IN
    GPIO --> DATA_IO
    UART --> CONSOLE
    USB_PORTS --> DUT

    RELAY --> POWER_IN
    RELAY --> DIP_SW
    SDMUX --> SD_SLOT
    SERIAL --> CONSOLE
    ETH_USB --> DUT

    MONITOR --> MCU

    style HOST fill:#e1f5fe
    style MCU fill:#fff3e0
    style PI fill:#e8f5e8
    style DUT fill:#ffebee
    style POWER fill:#e8f5e8
    style RELAY fill:#ffcdd2
```

## DUT Link Board

## Original DUT Link Board

The original DUT Link Board is a custom hardware solution for interfacing with devices under test.

### Features

- Multiple I/O interfaces
- Power control and monitoring
- Real-time data acquisition
- USB and Ethernet connectivity

### Specifications

#### Power

- Input: 12V DC
- Output: Configurable 3.3V, 5V, 12V
- Current monitoring: Up to 10A
- Protection: Over-current, over-voltage

#### I/O Interfaces

- GPIO: 16 configurable pins
- UART: 4 ports
- SPI: 2 ports
- I2C: 2 ports
- USB: 4 ports (host/device configurable)

#### Connectivity

- Ethernet: 1Gbps
- USB 3.0: Host connection
- Expansion headers for custom interfaces

### Board Layout

```mermaid
graph TB
    subgraph "DUT Link Board Physical Layout"
        subgraph "Top Section"
            PWR_IN[12V Power Input<br/>Barrel Jack]
            USB_HOST[USB 3.0 Host<br/>Connection]
        end

        subgraph "Center Section"
            MCU_AREA[STM32 MCU<br/>Main Controller]
            GPIO_HDR[GPIO Headers<br/>16 Pins]
        end

        subgraph "Left Section"
            UART_CON[UART Connectors<br/>4 Ports]
            SPI_CON[SPI Connectors<br/>2 Ports]
            I2C_CON[I2C Connectors<br/>2 Ports]
        end

        subgraph "Right Section"
            USB_PORTS[USB Ports<br/>4x Type-A]
            STATUS_LED[Status LEDs<br/>Power, Activity, Error]
        end

        subgraph "Bottom Section"
            ETH_PORT[Ethernet Port<br/>1Gbps RJ45]
            EXPANSION[Expansion Headers<br/>Custom Interfaces]
        end

        PWR_IN --> MCU_AREA
        USB_HOST --> MCU_AREA
        MCU_AREA --> GPIO_HDR
        MCU_AREA --> UART_CON
        MCU_AREA --> SPI_CON
        MCU_AREA --> I2C_CON
        MCU_AREA --> USB_PORTS
        MCU_AREA --> STATUS_LED
        MCU_AREA --> ETH_PORT
        MCU_AREA --> EXPANSION
    end

    style MCU_AREA fill:#fff3e0
    style PWR_IN fill:#e8f5e8
    style USB_HOST fill:#e1f5fe
    style STATUS_LED fill:#ffebee
```

## DUT Pi-Link Board

The DUT Pi-Link Board is a Raspberry Pi-based hardware solution that provides additional capabilities for device testing and automation, particularly focused on power management, DIP switch control, and SD card flashing.

### Architecture Overview

```mermaid
graph TB
    subgraph "DUT Pi-Link Board"
        subgraph "Raspberry Pi"
            PI_CPU[Raspberry Pi<br/>Linux System]
            PI_USB[USB Ports]
            PI_ETH[Ethernet Port<br/>External Network]
            PI_GPIO[40-Pin GPIO<br/>Header]
        end

        subgraph "Waveshare 8-Channel Relay"
            RELAY_1[Relay 1<br/>Power Switch]
            RELAY_2[Relay 2<br/>DIP Switch 1]
            RELAY_3[Relay 3<br/>DIP Switch 2]
            RELAY_4[Relay 4<br/>DIP Switch 3]
            RELAY_5[Relay 5<br/>DIP Switch 4]
            RELAY_6[Relay 6<br/>DIP Switch 5]
            RELAY_7[Relay 7<br/>DIP Switch 6]
            RELAY_8[Relay 8<br/>DIP Switch 7]
        end

        subgraph "Interface Components"
            SDMUX[USB-SD-Mux FAST<br/>SD Card Flashing]
            SERIAL[USB-Serial<br/>Console Converter]
            ETH_USB[Ethernet-USB<br/>Dongle]
        end
    end

    subgraph "Target Device"
        DUT_DEVICE[Device Under Test]
        DUT_POWER[Power Input]
        DUT_SD[SD Card Slot]
        DUT_CONSOLE[Serial Console]
        DUT_NET[Network Interface]
        DUT_DIP[DIP Switches<br/>Boot Configuration]
    end

    PI_CPU --> PI_GPIO
    PI_GPIO --> RELAY_1
    PI_GPIO --> RELAY_2
    PI_GPIO --> RELAY_3
    PI_GPIO --> RELAY_4
    PI_GPIO --> RELAY_5
    PI_GPIO --> RELAY_6
    PI_GPIO --> RELAY_7
    PI_GPIO --> RELAY_8

    PI_USB --> SDMUX
    PI_USB --> SERIAL
    PI_USB --> ETH_USB

    RELAY_1 --> DUT_POWER
    RELAY_2 --> DUT_DIP
    RELAY_3 --> DUT_DIP
    RELAY_4 --> DUT_DIP
    RELAY_5 --> DUT_DIP
    RELAY_6 --> DUT_DIP
    RELAY_7 --> DUT_DIP
    RELAY_8 --> DUT_DIP

    SDMUX --> DUT_SD
    SERIAL --> DUT_CONSOLE
    ETH_USB --> DUT_NET

    style PI_CPU fill:#e8f5e8
    style RELAY_1 fill:#ffcdd2
    style SDMUX fill:#e1f5fe
    style DUT_DEVICE fill:#ffebee
```

### Key Features

#### Power Management

- **Relay Channel 1**: Primary power switch capability
- **Clean power cycling**: Software-controlled power on/off sequences
- **Power sequencing**: Configurable delays for proper device initialization

#### DIP Switch Override

- **Relay Channels 2-8**: Override individual DIP switches
- **Boot configuration control**: Change device boot modes remotely
- **Hardware configuration**: Modify device settings without physical access
- **Test mode selection**: Switch between different operational modes

#### SD Card Management

- **USB-SD-Mux FAST**: High-speed SD card flashing and switching
- **Remote imaging**: Flash new firmware/OS images without physical access
- **Boot media control**: Switch between different boot images
- **Fast switching**: Quick transition between host and target SD access

#### Console Access

- **USB-Serial Converter**: Direct access to device debug console
- **Remote monitoring**: Capture boot logs and runtime messages
- **Interactive debugging**: Send commands to device console
- **Log collection**: Automated capture of device output

#### Network Configuration

- **Dual Ethernet Setup**:
  - **External Network**: Pi's built-in Ethernet for management
  - **DUT Internal Network**: USB-Ethernet dongle for device communication
- **Network isolation**: Separate management and test networks
- **DHCP/Static IP**: Flexible IP configuration for both networks

### Hardware Specifications

#### Raspberry Pi Base

- **Model**: Raspberry Pi 4B (4GB+ recommended)
- **Storage**: 32GB+ microSD card for Pi OS
- **Power**: 5V/3A USB-C power supply
- **Connectivity**: Built-in WiFi, Bluetooth, Ethernet

#### Waveshare 8-Channel Relay Board

- **Relay Type**: SPDT (Single Pole Double Throw)
- **Contact Rating**: 10A/250VAC, 10A/30VDC
- **Control Voltage**: 3.3V/5V compatible
- **Interface**: GPIO control via Pi header

#### USB-SD-Mux FAST Specifications

- **Switching Speed**: <1 second
- **USB Interface**: USB 3.0 SuperSpeed
- **SD Card Support**: SDHC/SDXC up to 2TB
- **Remote Control**: Software-controlled switching

### Software Integration

#### Jumpstarter Integration

```python
from jumpstarter.hardware import PiLinkBoard

# Initialize Pi-Link Board
board = PiLinkBoard(
    hostname="pi-link-001.local",
    relay_channels=8,
    sd_mux=True,
    serial_console=True
)

# Power cycle with DIP switch configuration
board.power_off()
board.set_dip_switches([1, 0, 1, 0, 1, 0, 0])  # Configure boot mode
board.flash_sd_card("firmware_v2.1.img")
board.power_on()

# Monitor console output
console_log = board.read_console(timeout=30)
print(f"Boot log: {console_log}")
```

#### Configuration Management

```yaml
# pi-link-config.yaml
board:
  hostname: "pi-link-001.local"

power:
  relay_channel: 1
  power_on_delay: 2.0
  power_off_delay: 1.0

dip_switches:
  relay_channels: [2, 3, 4, 5, 6, 7, 8]
  boot_modes:
    normal: [0, 0, 0, 0, 0, 0, 0]
    recovery: [1, 0, 0, 0, 0, 0, 0]
    factory: [0, 1, 0, 0, 0, 0, 0]
    test: [1, 1, 0, 0, 0, 0, 0]

sd_mux:
  device: "/dev/usbsdmux-001"
  mount_point: "/mnt/dut-sd"

console:
  device: "/dev/ttyUSB0"
  baudrate: 115200
  timeout: 30

network:
  management:
    interface: "eth0"
    ip: "192.168.1.100"
  dut:
    interface: "eth1" # USB-Ethernet adapter
    ip: "10.0.0.1"
    dhcp_range: "10.0.0.10-10.0.0.50"
```

### Setup and Configuration

#### Initial Setup

1. **Prepare Raspberry Pi**:

   ```bash
   # Flash Raspberry Pi OS
   sudo dd if=raspios-lite.img of=/dev/sdX bs=4M status=progress

   # Enable SSH and configure networking
   touch /boot/ssh
   echo 'pi:$encrypted_password' > /boot/userconf.txt
   ```

2. **Install Dependencies**:

   ```bash
   sudo apt update
   sudo apt install python3-pip python3-venv git
   sudo pip3 install jumpstarter-pi-link
   ```

3. **Configure Hardware**:
   ```bash
   # Enable GPIO and serial interfaces
   sudo raspi-config nonint do_spi 0
   sudo raspi-config nonint do_serial 0
   sudo raspi-config nonint do_ssh 0
   ```

#### Relay Board Connection

```bash
# GPIO Pin Mapping for Waveshare 8-Channel Relay
Relay 1 (Power): GPIO 26
Relay 2 (DIP 1): GPIO 20
Relay 3 (DIP 2): GPIO 21
Relay 4 (DIP 3): GPIO 22
Relay 5 (DIP 4): GPIO 23
Relay 6 (DIP 5): GPIO 24
Relay 7 (DIP 6): GPIO 25
Relay 8 (DIP 7): GPIO 27
```

### Use Cases

#### Automated Firmware Testing

```python
# Test multiple firmware versions
firmware_images = [
    "firmware_v1.0.img",
    "firmware_v1.1.img",
    "firmware_v2.0.img"
]

for image in firmware_images:
    board.power_off()
    board.flash_sd_card(image)
    board.set_dip_switches([0, 0, 0, 0, 0, 0, 0])  # Normal boot
    board.power_on()

    # Wait for boot and run tests
    if board.wait_for_boot(timeout=60):
        test_results = run_firmware_tests(board)
        log_results(image, test_results)
```

#### Boot Mode Testing

```python
# Test different boot configurations
boot_modes = {
    "normal": [0, 0, 0, 0, 0, 0, 0],
    "recovery": [1, 0, 0, 0, 0, 0, 0],
    "factory": [0, 1, 0, 0, 0, 0, 0]
}

for mode, dip_config in boot_modes.items():
    board.power_off()
    board.set_dip_switches(dip_config)
    board.power_on()

    boot_log = board.read_console(timeout=30)
    validate_boot_mode(mode, boot_log)
```

### Troubleshooting

#### Common Issues

1. **Relay Not Switching**
   - Check GPIO pin connections
   - Verify relay board power supply
   - Test individual GPIO pins

2. **SD-Mux Not Working**
   - Check USB connection
   - Verify device permissions
   - Update usbsdmux firmware

3. **Console Not Accessible**
   - Check USB-serial adapter connection
   - Verify device permissions (`sudo usermod -a -G dialout $USER`)
   - Test with different baud rates

4. **Network Issues**
   - Verify USB-Ethernet adapter driver
   - Check network configuration
   - Test connectivity with ping

#### Debug Commands

```bash
# Test relay control
gpio -g mode 26 out
gpio -g write 26 1  # Turn on relay 1
gpio -g write 26 0  # Turn off relay 1

# Check SD-Mux status
usbsdmux /dev/usbsdmux-001 get

# Test serial console
minicom -D /dev/ttyUSB0 -b 115200

# Check network interfaces
ip addr show
ping -c 3 10.0.0.10
```

### Safety and Best Practices

#### Electrical Safety

- Use appropriate relay ratings for your DUT power requirements
- Implement proper isolation between relay contacts and control circuits
- Add fuses/circuit breakers for overcurrent protection

#### Software Safety

- Implement timeout mechanisms for all operations
- Add retry logic for hardware communication
- Log all control operations for debugging

#### Maintenance

- Regular testing of relay switching cycles
- Periodic cleaning of relay contacts
- SD card wear monitoring and replacement
- Network connectivity verification

## Original DUT Link Board

## DUT Link Firmware

The firmware is written in Rust and provides low-level hardware control.

### Firmware Architecture

```mermaid
graph TB
    subgraph "Firmware Stack"
        subgraph "Application Layer"
            APP[Main Control Logic<br/>Command Processing]
            STATE[State Management<br/>Device Status]
        end

        subgraph "Protocol Layer"
            COMM[Communication Handler<br/>USB/Ethernet]
            CMD[Command Parser<br/>Message Routing]
            RESP[Response Generator<br/>Status Reporting]
        end

        subgraph "Driver Layer"
            GPIO_DRV[GPIO Driver]
            UART_DRV[UART Driver]
            SPI_DRV[SPI Driver]
            I2C_DRV[I2C Driver]
            USB_DRV[USB Driver]
            PWR_DRV[Power Driver]
        end

        subgraph "HAL Layer"
            HAL[Hardware Abstraction<br/>STM32 HAL]
            RTOS[Real-time Kernel<br/>Embedded Runtime]
        end
    end

    APP --> CMD
    STATE --> RESP
    CMD --> COMM
    RESP --> COMM

    CMD --> GPIO_DRV
    CMD --> UART_DRV
    CMD --> SPI_DRV
    CMD --> I2C_DRV
    CMD --> USB_DRV
    CMD --> PWR_DRV

    GPIO_DRV --> HAL
    UART_DRV --> HAL
    SPI_DRV --> HAL
    I2C_DRV --> HAL
    USB_DRV --> HAL
    PWR_DRV --> HAL

    HAL --> RTOS

    style APP fill:#e1f5fe
    style COMM fill:#fff3e0
    style HAL fill:#e8f5e8
    style RTOS fill:#ffebee
```

### Features

- Real-time operation
- Low latency communication
- Robust error handling
- Firmware update capability

### Communication Protocol

The firmware uses a custom protocol over USB/Ethernet:

```rust
#[derive(Debug, Serialize, Deserialize)]
pub enum Command {
    PowerOn { port: u8 },
    PowerOff { port: u8 },
    ReadGpio { pin: u8 },
    WriteGpio { pin: u8, value: bool },
    // ... more commands
}
```

## Setup and Configuration

### Initial Setup

1. Connect the DUT Link Board to your host system
2. Flash the firmware using the programming interface
3. Configure network settings if using Ethernet

### Firmware Updates

Update firmware using the built-in bootloader:

```bash
jumpstarter firmware update dutlink-board.bin
```

### Calibration

Calibrate voltage and current measurements:

```bash
jumpstarter calibrate --board dutlink-001
```

## Troubleshooting

### Common Issues

1. **Board Not Detected**
   - Check USB/Ethernet connections
   - Verify drivers are installed
   - Check power supply

2. **Communication Errors**
   - Verify firmware version compatibility
   - Check cable integrity
   - Review network configuration

3. **Power Issues**
   - Check input voltage (12V ±5%)
   - Verify current limits
   - Check for short circuits

### Debug Tools

- **Status LEDs**: Indicate board state
- **Debug UART**: Low-level debugging
- **Web Interface**: Configuration and monitoring

## Extension and Customization

### Custom Drivers

Add support for new devices:

```rust
impl DeviceDriver for CustomDevice {
    fn initialize(&mut self) -> Result<(), Error> {
        // Implementation
        Ok(())
    }

    fn reset(&mut self) -> Result<(), Error> {
        // Implementation
        Ok(())
    }
}
```

### Expansion Boards

Create expansion boards for specialized testing:

- Analog test interfaces
- High-speed digital interfaces
- RF test capabilities
- Environmental sensors

## Safety Considerations

- Always verify connections before powering on
- Use appropriate current limits
- Follow ESD protection procedures
- Ensure proper grounding
