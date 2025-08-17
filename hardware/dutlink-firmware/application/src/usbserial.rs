use stm32f4xx_hal::nb;
use stm32f4xx_hal::otg_fs::UsbBusType;
use usb_device::{class_prelude::*, LangID, Result};
use usbd_serial::SerialPort;

// Bigger USB Serial buffer
use core::borrow::{Borrow, BorrowMut};

pub const BUFFER_SIZE: usize = 1024;
pub struct BufferStore(pub [u8; BUFFER_SIZE]);

impl Borrow<[u8]> for BufferStore {
    fn borrow(&self) -> &[u8] {
        &self.0
    }
}

impl BorrowMut<[u8]> for BufferStore {
    fn borrow_mut(&mut self) -> &mut [u8] {
        &mut self.0
    }
}

// Wrapper for SerialPort that implements ushell traits
pub struct UShellSerialWrapper {
    serial: SerialPort<'static, UsbBusType, BufferStore, BufferStore>,
}

impl UShellSerialWrapper {
    pub fn new(serial: SerialPort<'static, UsbBusType, BufferStore, BufferStore>) -> Self {
        Self { serial }
    }

    // Delegate methods to the underlying SerialPort
    pub fn write(&mut self, data: &[u8]) -> usb_device::Result<usize> {
        self.serial.write(data)
    }

    pub fn read(&mut self, data: &mut [u8]) -> usb_device::Result<usize> {
        self.serial.read(data)
    }

    pub fn flush(&mut self) -> usb_device::Result<()> {
        self.serial.flush()
    }

    pub fn reset(&mut self) {
        self.serial.reset()
    }
}

// Implement UsbClass for the wrapper to delegate to the underlying SerialPort
impl UsbClass<UsbBusType> for UShellSerialWrapper {
    fn get_configuration_descriptors(&self, writer: &mut DescriptorWriter) -> Result<()> {
        self.serial.get_configuration_descriptors(writer)
    }

    fn get_bos_descriptors(&self, writer: &mut BosWriter) -> Result<()> {
        self.serial.get_bos_descriptors(writer)
    }

    fn get_string(&self, index: StringIndex, lang_id: LangID) -> Option<&str> {
        self.serial.get_string(index, lang_id)
    }

    fn reset(&mut self) {
        self.serial.reset()
    }

    fn set_alt_setting(&mut self, interface: InterfaceNumber, alt_setting: u8) -> bool {
        self.serial.set_alt_setting(interface, alt_setting)
    }

    fn control_in(&mut self, xfer: ControlIn<UsbBusType>) {
        self.serial.control_in(xfer)
    }

    fn control_out(&mut self, xfer: ControlOut<UsbBusType>) {
        self.serial.control_out(xfer)
    }

    fn endpoint_setup(&mut self, addr: EndpointAddress) {
        self.serial.endpoint_setup(addr)
    }

    fn endpoint_out(&mut self, addr: EndpointAddress) {
        self.serial.endpoint_out(addr)
    }

    fn endpoint_in_complete(&mut self, addr: EndpointAddress) {
        self.serial.endpoint_in_complete(addr)
    }
}

// Implement ushell::Read trait for the wrapper
impl ushell::Read<u8> for UShellSerialWrapper {
    type Error = ();

    fn read(&mut self) -> nb::Result<u8, Self::Error> {
        let mut buf = [0u8; 1];
        match self.serial.read(&mut buf) {
            Ok(1) => Ok(buf[0]),
            Ok(0) => Err(nb::Error::WouldBlock),
            Ok(_) => unreachable!(),
            Err(_) => Err(nb::Error::Other(())),
        }
    }
}

// Implement ushell::Write trait for the wrapper
impl ushell::Write<u8> for UShellSerialWrapper {
    type Error = ();

    fn write(&mut self, word: u8) -> nb::Result<(), Self::Error> {
        let buf = [word];
        match self.serial.write(&buf) {
            Ok(1) => Ok(()),
            Ok(0) => Err(nb::Error::WouldBlock),
            Ok(_) => unreachable!(),
            Err(_) => Err(nb::Error::Other(())),
        }
    }

    fn flush(&mut self) -> nb::Result<(), Self::Error> {
        match self.serial.flush() {
            Ok(()) => Ok(()),
            Err(_) => Err(nb::Error::Other(())),
        }
    }
}

pub type USBSerialType = UShellSerialWrapper;

macro_rules! new_usb_serial {
    ($usb:expr) => {
        UShellSerialWrapper::new(SerialPort::new_with_store(
            $usb,
            BufferStore([0; BUFFER_SIZE]),
            BufferStore([0; BUFFER_SIZE]),
        ))
    };
}
pub(crate) use new_usb_serial;
