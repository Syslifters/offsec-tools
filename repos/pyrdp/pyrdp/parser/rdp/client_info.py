#
# This file is part of the PyRDP project.
# Copyright (C) 2018, 2021 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

from io import BytesIO

from pyrdp.core import decodeUTF16LE, encodeUTF16LE, StrictStream, Uint16LE, Uint32LE
from pyrdp.enum import ClientInfoFlags
from pyrdp.parser.parser import Parser
from pyrdp.pdu import ClientExtraInfo, ClientInfoPDU


class ClientInfoParser(Parser):
    """
    Read and write the RDP ClientInfo PDU which contains very useful information.
    See https://msdn.microsoft.com/en-us/library/cc240475.aspx
    """

    def doParse(self, data: bytes) -> ClientInfoPDU:
        """
        Decode a Client Info PDU from bytes.
        :param data: the Client Info PDU bytes.
        """
        stream = BytesIO(data)
        codePage = Uint32LE.unpack(stream)
        flags = Uint32LE.unpack(stream)

        isUnicode = flags & ClientInfoFlags.INFO_UNICODE != 0
        hasNullBytes = codePage == 1252 or isUnicode
        nullByteCount = 1 if hasNullBytes else 0
        unicodeMultiplier = 2 if isUnicode else 0

        domainLength = Uint16LE.unpack(stream) + nullByteCount * unicodeMultiplier
        usernameLength = Uint16LE.unpack(stream) + nullByteCount * unicodeMultiplier
        passwordLength = Uint16LE.unpack(stream) + nullByteCount * unicodeMultiplier
        alternateShellLength = Uint16LE.unpack(stream) + nullByteCount * unicodeMultiplier
        workingDirLength = Uint16LE.unpack(stream) + nullByteCount * unicodeMultiplier

        domain = decodeUTF16LE(stream.read(domainLength))
        username = decodeUTF16LE(stream.read(usernameLength))
        password = decodeUTF16LE(stream.read(passwordLength))
        alternateShell = decodeUTF16LE(stream.read(alternateShellLength))
        workingDir = decodeUTF16LE(stream.read(workingDirLength))

        extraInfoBytes = stream.read()

        if extraInfoBytes != b"":
            extraInfo = self.parseExtraInfo(extraInfoBytes)
        else:
            extraInfo = None

        return ClientInfoPDU(codePage, flags, domain, username, password, alternateShell, workingDir, extraInfo)

    def write(self, pdu: ClientInfoPDU) -> bytes:
        """
        Encode a Client Info PDU to bytes.
        :param pdu: the Client Info PDU.
        """
        stream = BytesIO()
        stream.write(Uint32LE.pack(pdu.codePage))
        stream.write(Uint32LE.pack(pdu.flags))

        isUnicode = pdu.flags & ClientInfoFlags.INFO_UNICODE != 0
        hasNullBytes = pdu.codePage == 1252 or isUnicode
        nullByteCount = 1 if hasNullBytes else 0
        unicodeMultiplier = 2 if isUnicode else 0

        domain = pdu.domain + "\x00" * nullByteCount
        username = pdu.username + "\x00" * nullByteCount
        password = pdu.password + "\x00" * nullByteCount
        alternateShell = pdu.alternateShell + "\x00" * nullByteCount
        workingDir = pdu.workingDir + "\x00" * nullByteCount

        if isUnicode:
            domain = encodeUTF16LE(domain)
            username = encodeUTF16LE(username)
            password = encodeUTF16LE(password)
            alternateShell = encodeUTF16LE(alternateShell)
            workingDir = encodeUTF16LE(workingDir)
        else:
            domain = domain.encode()
            username = username.encode()
            password = password.encode()
            alternateShell = alternateShell.encode()
            workingDir = workingDir.encode()

        domainLength = len(domain) - nullByteCount * unicodeMultiplier
        usernameLength = len(username) - nullByteCount * unicodeMultiplier
        passwordLength = len(password) - nullByteCount * unicodeMultiplier
        alternateShellLength = len(alternateShell) - nullByteCount * unicodeMultiplier
        workingDirLength = len(workingDir) - nullByteCount * unicodeMultiplier

        stream.write(Uint16LE.pack(domainLength))
        stream.write(Uint16LE.pack(usernameLength))
        stream.write(Uint16LE.pack(passwordLength))
        stream.write(Uint16LE.pack(alternateShellLength))
        stream.write(Uint16LE.pack(workingDirLength))
        stream.write(domain)
        stream.write(username)
        stream.write(password)
        stream.write(alternateShell)
        stream.write(workingDir)

        if pdu.extraInfo is not None:
            extraInfoBytes = self.writeExtraInfo(pdu.extraInfo)
            stream.write(extraInfoBytes)

        return stream.getvalue()

    def parseExtraInfo(self, data: bytes) -> ClientExtraInfo:
        stream = BytesIO(data)
        clientAddressFamily = Uint16LE.unpack(stream)

        clientAddressLength = Uint16LE.unpack(stream)
        clientAddress = stream.read(clientAddressLength)

        clientDirLength = Uint16LE.unpack(stream)
        clientDir = stream.read(clientDirLength)

        extraInfo = ClientExtraInfo(clientAddressFamily, clientAddress, clientDir)
        stream = StrictStream(stream)

        try:
            extraInfo.clientTimeZone = stream.read(172)
            extraInfo.clientSessionID = Uint32LE.unpack(stream)
            extraInfo.performanceFlags = Uint32LE.unpack(stream)

            autoReconnectCookieLength = Uint16LE.unpack(stream)
            extraInfo.autoReconnectCookie = stream.read(autoReconnectCookieLength)

            stream.read(4)

            dynamicDSTTimeZoneKeyNameLength = Uint16LE.unpack(stream)
            extraInfo.dynamicDSTTimeZoneKeyName = stream.read(dynamicDSTTimeZoneKeyNameLength)

            extraInfo.dynamicDaylightTimeDisabled = bool(Uint16LE.unpack(stream))
        except EOFError:
            pass

        return extraInfo

    def writeExtraInfo(self, extraInfo: ClientExtraInfo) -> bytes:
        stream = BytesIO()
        Uint16LE.pack(extraInfo.clientAddressFamily, stream)

        Uint16LE.pack(len(extraInfo.clientAddress), stream)
        stream.write(extraInfo.clientAddress)

        Uint16LE.pack(len(extraInfo.clientDir), stream)
        stream.write(extraInfo.clientDir)

        if extraInfo.clientTimeZone is None:
            return stream.getvalue()

        stream.write(extraInfo.clientTimeZone)

        if extraInfo.clientSessionID is None:
            return stream.getvalue()

        Uint32LE.pack(extraInfo.clientSessionID, stream)

        if extraInfo.performanceFlags is None:
            return stream.getvalue()

        Uint32LE.pack(extraInfo.performanceFlags, stream)

        if extraInfo.autoReconnectCookie is None:
            return stream.getvalue()

        Uint16LE.pack(len(extraInfo.autoReconnectCookie), stream)
        stream.write(extraInfo.autoReconnectCookie)

        stream.write(b"\x00" * 4)

        if extraInfo.dynamicDSTTimeZoneKeyName is None:
            return stream.getvalue()

        Uint16LE.pack(len(extraInfo.dynamicDSTTimeZoneKeyName), stream)
        stream.write(extraInfo.dynamicDSTTimeZoneKeyName)

        if extraInfo.dynamicDaylightTimeDisabled is None:
            return stream.getvalue()

        Uint16LE.pack(int(extraInfo.dynamicDaylightTimeDisabled), stream)
        return stream.getvalue()