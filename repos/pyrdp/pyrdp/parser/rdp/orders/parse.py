#
# This file is part of the PyRDP project.
# Copyright (C) 2020 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

"""
Parse Drawing Orders.
"""
import logging
from io import BytesIO

from pyrdp.core.packing import Uint8, Uint16LE
from pyrdp.pdu.rdp.fastpath import FastPathOrdersEvent
from pyrdp.enum.orders import DrawingOrderControlFlags as ControlFlags
from pyrdp.enum.rdp import GlyphSupport
from pyrdp.pdu.rdp.capability import CapabilityType

from .frontend import GdiFrontend

from .secondary import CacheBitmapV1, CacheColorTable, CacheGlyph, CacheBitmapV2, CacheBrush, CacheBitmapV3
from .alternate import CreateOffscreenBitmap, SwitchSurface, CreateNineGridBitmap, \
    StreamBitmapFirst, StreamBitmapNext, GdiPlusFirst, GdiPlusNext, GdiPlusEnd, GdiPlusCacheFirst, \
    GdiPlusCacheNext, GdiPlusCacheEnd, FrameMarker
from .primary import PrimaryContext as Context

LOG = logging.getLogger(__name__)


def _repr(n):
    """Internal method to stringify an order type."""
    r = n.__doc__
    return r if r else 'UNKNOWN (%02x)'.format(n)


class OrdersParser:
    """
    Drawing Order Parser.
    """

    def __init__(self, frontend: GdiFrontend):
        """
        Create a drawing order parser.

        :param GdiFrontend frontend: The frontend that will process GDI messages.
        """

        self.notify: GdiFrontend = frontend
        self.ctx = Context()
        self.glyphLevel: GlyphSupport = GlyphSupport.GLYPH_SUPPORT_NONE

    def onCapabilities(self, caps):
        """Update the parser to take into account the capabilities reported by the client."""

        if CapabilityType.CAPSTYPE_GLYPHCACHE in caps:
            glyphLevel = caps[CapabilityType.CAPSTYPE_GLYPHCACHE].glyphSupportLevel
            self.glyphLevel = GlyphSupport(glyphLevel)

    def parse(self, orders: FastPathOrdersEvent):
        """
        Entrypoint for parsing TS_FP_UPDATE_ORDERS.
        """

        s = BytesIO(orders.payload)

        numberOrders = Uint16LE.unpack(s)
        try:
            for _ in range(numberOrders):
                self._parse_order(s)
        except Exception:
            LOG.warning('Failed to parse drawing order PDU: %s', orders)

        return orders

    def _parse_order(self, s: BytesIO):
        controlFlags = Uint8.unpack(s)

        if not (controlFlags & ControlFlags.TS_STANDARD):
            self._parse_altsec(s, controlFlags)
        elif (controlFlags & ControlFlags.TS_SECONDARY):
            self._parse_secondary(s, controlFlags)
        else:
            self._parse_primary(s, controlFlags)

    # Primary drawing orders.
    # ----------------------------------------------------------------------
    def _parse_primary(self, s: BytesIO, flags: int):

        orderType = self.ctx.update(s, flags)
        self.notify.onBounds(self.ctx.bounds if self.ctx.bounded else None)

        assert orderType >= 0 and orderType < len(_pri)
        _pri[orderType](self, s)

    def _parse_dstblt(self, s: BytesIO):
        """DSTBLT"""
        self.notify.dstBlt(self.ctx.dstBlt.update(s))

    def _parse_patblt(self, s: BytesIO):
        """PATBLT"""
        self.notify.patBlt(self.ctx.patBlt.update(s))

    def _parse_scrblt(self, s: BytesIO):
        """SCRBLT"""
        self.notify.scrBlt(self.ctx.scrBlt.update(s))

    def _parse_draw_nine_grid(self, s: BytesIO):
        """DRAW_NINE_GRID"""
        self.notify.drawNineGrid(self.ctx.drawNineGrid.update(s))

    def _parse_multi_draw_nine_grid(self, s: BytesIO):
        """MULTI_DRAW_NINE_GRID"""
        self.notify.multiDrawNineGrid(self.ctx.multiDrawNineGrid.update(s))

    def _parse_line_to(self, s: BytesIO):
        """LINE_TO"""
        self.notify.lineTo(self.ctx.lineTo.update(s))

    def _parse_opaque_rect(self, s: BytesIO):
        """OPAQUE_RECT"""
        self.notify.opaqueRect(self.ctx.opaqueRect.update(s))

    def _parse_save_bitmap(self, s: BytesIO):
        """SAVE_BITMAP"""
        self.notify.saveBitmap(self.ctx.saveBitmap.update(s))

    def _parse_memblt(self, s: BytesIO):
        """MEMBLT"""
        self.notify.memBlt(self.ctx.memBlt.update(s))

    def _parse_mem3blt(self, s: BytesIO):
        """MEM3BLT"""
        self.notify.mem3Blt(self.ctx.mem3Blt.update(s))

    def _parse_multi_dstblt(self, s: BytesIO):
        """MULTI_DSTBLT"""
        self.notify.multiDstBlt(self.ctx.multiDstBlt.update(s))

    def _parse_multi_patblt(self, s: BytesIO):
        """MULTI_PATBLT"""
        self.notify.multiPatBlt(self.ctx.multiPatBlt.update(s))

    def _parse_multi_scrblt(self, s: BytesIO):
        """MULTI_SCRBLT"""
        self.notify.multiScrBlt(self.ctx.multiScrBlt.update(s))

    def _parse_multi_opaque_rect(self, s: BytesIO):
        """MULTI_OPAQUE_RECT"""
        self.notify.multiOpaqueRect(self.ctx.multiOpaqueRect.update(s))

    def _parse_fast_index(self, s: BytesIO):
        """FAST_INDEX"""
        self.notify.fastIndex(self.ctx.fastIndex.update(s))

    def _parse_polygon_sc(self, s: BytesIO):
        """POLYGON_SC"""
        self.notify.polygonSc(self.ctx.polygonSc.update(s))

    def _parse_polygon_cb(self, s: BytesIO):
        """POLYGON_CB"""
        self.notify.polygonCb(self.ctx.polygonCb.update(s))

    def _parse_polyLine(self, s: BytesIO):
        """POLYLINE"""
        self.notify.polyLine(self.ctx.polyLine.update(s))

    def _parse_fast_glyph(self, s: BytesIO):
        """FAST_GLYPH"""
        self.notify.fastGlyph(self.ctx.fastGlyph.update(s))

    def _parse_ellipse_sc(self, s: BytesIO):
        """ELLIPSE_SC"""
        self.notify.ellipseSc(self.ctx.ellipseSc.update(s))

    def _parse_ellipse_cb(self, s: BytesIO):
        """ELLIPSE_CB"""
        self.notify.ellipseCb(self.ctx.ellipseCb.update(s))

    def _parse_glyph_index(self, s: BytesIO):
        """GLYPH_INDEX"""
        self.notify.glyphIndex(self.ctx.glyphIndex.update(s))

    # Secondary drawing orders.
    # ----------------------------------------------------------------------
    def _parse_secondary(self, s: BytesIO, flags: int):
        Uint16LE.unpack(s)  # orderLength (unused)
        extraFlags = Uint16LE.unpack(s)
        orderType = Uint8.unpack(s)

        assert orderType >= 0 and orderType < len(_sec)
        _sec[orderType](self, s, orderType, extraFlags)

    def _parse_cache_bitmap_v1(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_BITMAP_V1"""
        self.notify.cacheBitmapV1(CacheBitmapV1.parse(s, orderType, flags))

    def _parse_cache_color_table(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_COLOR_TABLE"""
        self.notify.cacheColorTable(CacheColorTable.parse(s, orderType, flags))

    def _parse_cache_glyph(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_GLYPH"""
        if self.glyphLevel == GlyphSupport.GLYPH_SUPPORT_NONE:
            LOG.warn("Received CACHE_GLYPH but the client reported it doesn't support it!")
            # Ignore it.
        else:
            self.notify.cacheGlyph(CacheGlyph.parse(s, flags, self.glyphLevel))

    def _parse_cache_bitmap_v2(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_BITMAP_V2"""
        self.notify.cacheBitmapV2(CacheBitmapV2.parse(s, orderType, flags))

    def _parse_cache_brush(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_BRUSH"""
        self.notify.cacheBrush(CacheBrush.parse(s))

    def _parse_cache_bitmap_v3(self, s: BytesIO, orderType: int, flags: int):
        """CACHE_BITMAP_V3"""
        self.notify.cacheBitmapV3(CacheBitmapV3.parse(s, flags))

    # Alternate secondary drawing orders.
    # ----------------------------------------------------------------------
    def _parse_altsec(self, s: BytesIO, flags: int):
        orderType = flags >> 2

        assert orderType >= 0 and orderType < len(_alt)

        _alt[orderType](self, s)

    def _parse_create_offscreen_bitmap(self, s: BytesIO):
        """CREATE_OFFSCREEN_BITMAP"""
        self.notify.createOffscreenBitmap(CreateOffscreenBitmap.parse(s))

    def _parse_switch_surface(self, s: BytesIO):
        """SWITCH_SURFACE"""
        self.notify.switchSurface(SwitchSurface.parse(s))

    def _parse_create_nine_grid_bitmap(self, s: BytesIO):
        """CREATE_NINEGRID_BITMAP"""
        self.notify.createNineGridBitmap(CreateNineGridBitmap.parse(s))

    def _parse_stream_bitmap_first(self, s: BytesIO):
        """STREAM_BITMAP_FIRST"""
        self.notify.streamBitmapFirst(StreamBitmapFirst.parse(s))

    def _parse_stream_bitmap_next(self, s: BytesIO):
        """STREAM_BITMAP_NEXT"""
        self.notify.streamBitmapNext(StreamBitmapNext.parse(s))

    def _parse_gdiplus_first(self, s: BytesIO):
        """GDIPLUS_FIRST"""
        self.notify.gdiPlusFirst(GdiPlusFirst.parse(s))

    def _parse_gdiplus_next(self, s: BytesIO):
        """GDIPLUS_NEXT"""
        self.notify.drawGdiPlusNext(GdiPlusNext.parse(s))

    def _parse_gdiplus_end(self, s: BytesIO):
        """GDIPLUS_END"""
        self.notify.drawGdiPlusEnd(GdiPlusEnd.parse(s))

    def _parse_gdiplus_cache_first(self, s: BytesIO):
        """GDIPLUS_CACHE_FIRST"""
        self.notify.drawGdiPlusCacheFirst(GdiPlusCacheFirst.parse(s))

    def _parse_gdiplus_cache_next(self, s: BytesIO):
        """GDIPLUS_CACHE_NEXT"""
        self.notify.drawGdiPlusCacheNext(GdiPlusCacheNext.parse(s))

    def _parse_gdiplus_cache_end(self, s: BytesIO):
        """GDIPLUS_CACHE_END"""
        self.notify.drawGdiPlusCacheEnd(GdiPlusCacheEnd.parse(s))

    def _parse_window(self, s: BytesIO):
        """WINDOW"""
        # This is specified in MS-RDPERP for seamless applications.
        LOG.debug('WINDOW is not supported yet.')

    def _parse_compdesk_first(self, s: BytesIO):
        """COMPDESK"""
        LOG.debug('COMPDESK is not supported yet.')

    def _parse_frame_marker(self, s: BytesIO):
        """FRAME_MARKER"""
        self.notify.frameMarker(FrameMarker.parse(s))


# Parser Lookup Tables
_pri = [
    OrdersParser._parse_dstblt,                # 0x00
    OrdersParser._parse_patblt,                # 0x01
    OrdersParser._parse_scrblt,                # 0x02
    None,                                      # 0x03
    None,                                      # 0x04
    None,                                      # 0x05
    None,                                      # 0x06
    OrdersParser._parse_draw_nine_grid,        # 0x07
    OrdersParser._parse_multi_draw_nine_grid,  # 0x08
    OrdersParser._parse_line_to,               # 0x09
    OrdersParser._parse_opaque_rect,           # 0x0A
    OrdersParser._parse_save_bitmap,           # 0x0B
    None,                                      # 0x0C
    OrdersParser._parse_memblt,                # 0x0D
    OrdersParser._parse_mem3blt,               # 0x0E
    OrdersParser._parse_multi_dstblt,          # 0x0F
    OrdersParser._parse_multi_patblt,          # 0x10
    OrdersParser._parse_multi_scrblt,          # 0x11
    OrdersParser._parse_multi_opaque_rect,     # 0x12
    OrdersParser._parse_fast_index,            # 0x13
    OrdersParser._parse_polygon_sc,            # 0x14
    OrdersParser._parse_polygon_cb,            # 0x15
    OrdersParser._parse_polyLine,              # 0x16
    None,                                      # 0x17
    OrdersParser._parse_fast_glyph,            # 0x18
    OrdersParser._parse_ellipse_sc,            # 0x19
    OrdersParser._parse_ellipse_cb,            # 0x1A
    OrdersParser._parse_glyph_index,           # 0x1B
]

_sec = [
    OrdersParser._parse_cache_bitmap_v1,    # 0x00 : Uncompressed
    OrdersParser._parse_cache_color_table,  # 0x01
    OrdersParser._parse_cache_bitmap_v1,    # 0x02 : Compressed
    OrdersParser._parse_cache_glyph,        # 0x03
    OrdersParser._parse_cache_bitmap_v2,    # 0x04 : Uncompresed
    OrdersParser._parse_cache_bitmap_v2,    # 0x05 : Compressed
    None,                                   # 0x06
    OrdersParser._parse_cache_brush,        # 0x07
    OrdersParser._parse_cache_bitmap_v3,    # 0x08
]

_alt = [
    OrdersParser._parse_switch_surface,           # 0x00
    OrdersParser._parse_create_offscreen_bitmap,  # 0x01
    OrdersParser._parse_stream_bitmap_first,      # 0x02
    OrdersParser._parse_stream_bitmap_next,       # 0x03
    OrdersParser._parse_create_nine_grid_bitmap,  # 0x04
    OrdersParser._parse_gdiplus_first,            # 0x05
    OrdersParser._parse_gdiplus_next,             # 0x06
    OrdersParser._parse_gdiplus_end,              # 0x07
    OrdersParser._parse_gdiplus_cache_first,      # 0x08
    OrdersParser._parse_gdiplus_cache_next,       # 0x09
    OrdersParser._parse_gdiplus_cache_end,        # 0x0A
    OrdersParser._parse_window,                   # 0x0B
    OrdersParser._parse_compdesk_first,           # 0x0C
    OrdersParser._parse_frame_marker,             # 0x0D
]
