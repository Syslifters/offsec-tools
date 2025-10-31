#
# This file is part of the PyRDP project.
# Copyright (C) 2018 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

import logging
from binascii import hexlify

from pyrdp.security import SecuritySettingsObserver


class RC4LoggingObserver(SecuritySettingsObserver):
    def __init__(self, log: logging.LoggerAdapter):
        super().__init__()
        self.log = log

    def onCrypterGenerated(self, settings):
        self.log.debug("RC4 client/server random: %(rc4ClientRandom)s %(rc4ServerRandom)s",
                       {"rc4ClientRandom": hexlify(settings.clientRandom).decode(),
                        "rc4ServerRandom": hexlify(settings.serverRandom).decode()})