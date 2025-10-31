# This file is part of the PyRDP project.
#
# Copyright (C) 2020-2021 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

from pyrdp.enum import BitmapFlags, CapabilityType, SlowPathUpdateType
from pyrdp.parser import BitmapParser, FastPathOutputParser, OrdersParser
from pyrdp.pdu import BitmapUpdateData, FastPathBitmapEvent, FastPathOutputEvent, FastPathOrdersEvent, UpdatePDU
from pyrdp.player import ImageHandler
from pyrdp.player.BaseEventHandler import BaseEventHandler
from pyrdp.player.gdi.draw import GdiQtFrontend
from pyrdp.ui import RDPBitmapToQtImage

import logging


class RenderingEventHandler(BaseEventHandler):
    """Abstract class for rendering handlers."""

    def __init__(self, imageHandler: ImageHandler):
        BaseEventHandler.__init__(self)
        self._fastPath = FastPathOutputParser()
        self._bitmap = BitmapParser()
        self.gdi: GdiQtFrontend = None
        self._orders: OrdersParser = None
        self.log = logging.getLogger(__name__)
        self.imageHandler = imageHandler

    def onCapabilities(self, caps):
        if CapabilityType.CAPSTYPE_ORDER in caps:
            self.gdi = GdiQtFrontend(self.imageHandler)
            self._orders = OrdersParser(self.gdi)
            self._orders.onCapabilities(caps)

    # Generic Parsing Routines.
    def onFastPathOutput(self, event: FastPathOutputEvent):
        self.onBeginRender()
        if isinstance(event, FastPathBitmapEvent):
            parsed = self._fastPath.parseBitmapEvent(event)
            for bmp in parsed.bitmapUpdateData:
                self.onBitmap(bmp)

        elif isinstance(event, FastPathOrdersEvent):
            if self._orders is None:
                self.log.error('Received Unexpected Drawing Orders!')
                return
            self.onBeginRender()
            self._orders.parse(event)

        self.onFinishRender()

    def onSlowPathUpdate(self, pdu: UpdatePDU):
        if pdu.updateType == SlowPathUpdateType.SLOWPATH_UPDATETYPE_BITMAP:
            self.onBeginRender()

            updates = self._bitmap.parseBitmapUpdateData(pdu.updateData)
            for bmp in updates:
                self.onBitmap(bmp)

            self.onFinishRender()

    def onBitmap(self, bitmapData: BitmapUpdateData):
        image, _ = RDPBitmapToQtImage(
            bitmapData.width,
            bitmapData.heigth,
            bitmapData.bitsPerPixel,
            bitmapData.flags & BitmapFlags.BITMAP_COMPRESSION != 0,
            bitmapData.bitmapData
        )

        self.imageHandler.notifyImage(
            bitmapData.destLeft,
            bitmapData.destTop,
            image,
            bitmapData.destRight - bitmapData.destLeft + 1,
            bitmapData.destBottom - bitmapData.destTop + 1)

    def onBeginRender(self):
        pass

    def onFinishRender(self):
        pass
