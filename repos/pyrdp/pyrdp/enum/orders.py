#
# This file is part of the PyRDP project.
# Copyright (C) 2020 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

"""
Enumerations for Drawing Orders.
"""

from enum import IntEnum


class DrawingOrderControlFlags(IntEnum):
    """
    https://msdn.microsoft.com/en-us/library/cc241574.aspx
    """
    TS_STANDARD = 0x01
    TS_SECONDARY = 0x02
    TS_BOUNDS = 0x04
    TS_TYPE_CHANGE = 0x08
    TS_DELTA_COORDS = 0x10
    TS_ZERO_BOUNDS_DELTAS = 0x20
    TS_ZERO_FIELD_BYTE_BIT0 = 0x40
    TS_ZERO_FIELD_BYTE_BIT1 = 0x80


class Primary(IntEnum):
    DSTBLT = 0x00
    PATBLT = 0x01
    SCRBLT = 0x02
    DRAW_NINE_GRID = 0x07
    MULTI_DRAW_NINE_GRID = 0x08
    LINE_TO = 0x09
    OPAQUE_RECT = 0x0A
    SAVE_BITMAP = 0x0B
    MEMBLT = 0x0D
    MEM3BLT = 0x0E
    MULTI_DSTBLT = 0x0F
    MULTI_PATBLT = 0x10
    MULTI_SCRBLT = 0x11
    MULTI_OPAQUE_RECT = 0x12
    FAST_INDEX = 0x13
    POLYGON_SC = 0x14
    POLYGON_CB = 0x15
    POLYLINE = 0x16
    FAST_GLYPH = 0x18
    ELLIPSE_SC = 0x19
    ELLIPSE_CB = 0x1A
    GLYPH_INDEX = 0x1B


class Secondary(IntEnum):
    BITMAP_UNCOMPRESSED = 0x00
    CACHE_COLOR_TABLE = 0x01
    CACHE_BITMAP_COMPRESSED = 0x02
    CACHE_GLYPH = 0x03
    BITMAP_UNCOMPRESSED_V2 = 0x04
    BITMAP_COMPRESSED_V2 = 0x05
    CACHE_BRUSH = 0x07
    BITMAP_COMPRESSED_V3 = 0x08


class Alternate(IntEnum):
    SWITCH_SURFACE = 0x00
    CREATE_OFFSCREEN_BITMAP = 0x01
    STREAM_BITMAP_FIRST = 0x02
    STREAM_BITMAP_NEXT = 0x03
    CREATE_NINE_GRID_BITMAP = 0x04
    GDIPLUS_FIRST = 0x05
    GDIPLUS_NEXT = 0x06
    GDIPLUS_END = 0x07
    GDIPLUS_CACHE_FIRST = 0x08
    GDIPLUS_CACHE_NEXT = 0x09
    GDIPLUS_CACHE_END = 0x0A
    WINDOW = 0x0B
    COMPDESK_FIRST = 0x0C
    FRAME_MARKER = 0x0D
