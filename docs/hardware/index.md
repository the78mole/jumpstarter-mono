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
    
    subgraph "Device Under Test"
        DUT[Target Device]
        POWER_IN[Power Input]
        DATA_IO[Data I/O]
        CONSOLE[Debug Console]
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
    
    POWER --> POWER_IN
    GPIO --> DATA_IO
    UART --> CONSOLE
    USB_PORTS --> DUT
    
    MONITOR --> MCU
    
    style HOST fill:#e1f5fe
    style MCU fill:#fff3e0
    style DUT fill:#ffebee
    style POWER fill:#e8f5e8
```

## DUT Link Board

The DUT Link Board is a custom hardware solution for interfacing with devices under test.

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