#
# This file is part of the PyRDP project.
# Copyright (C) 2018 GoSecure Inc.
# Licensed under the GPLv3 or later.
#

from enum import IntEnum


class X224PDUType(IntEnum):
    """
    X224 header codes.
    """

    X224_TPDU_CONNECTION_REQUEST = 0x0E
    X224_TPDU_CONNECTION_CONFIRM = 0x0D
    X224_TPDU_DISCONNECT_REQUEST = 0x08
    X224_TPDU_DATA = 0x0F
    X224_TPDU_ERROR = 0x07
