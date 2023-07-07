﻿/*
 * BSD 3-Clause License
 *
 * Copyright (c) 2022, Kevin Robertson
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 * 
 * 1. Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its
 * contributors may be used to endorse or promote products derived from
 * this software without specific prior written permission. 
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Text;

namespace Quiddity.SMB
{
    class SMBCOMSessionSetupAndXRequest
    {
        // https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-smb/a00d0361-3544-4845-96ab-309b4bb7705d
        public byte WordCount { get; set; }
        public byte AndXCommand { get; set; }
        public byte AndXReserved { get; set; }
        public ushort AndXOffset { get; set; }
        public ushort MaxBufferSize { get; set; }
        public ushort MaxMpxCount { get; set; }
        public ushort VcNumber { get; set; }
        public uint SessionKey { get; set; }
        public ushort SecurityBlobLength { get; set; }
        public uint Reserved { get; set; }
        public uint Capabilities { get; set; }
        public ushort ByteCount { get; set; }
        public byte[] SecurityBlob { get; set; }

        public SMBCOMSessionSetupAndXRequest()
        {

        }

        public SMBCOMSessionSetupAndXRequest(byte[] data, int offset)
        {
            ReadBytes(data, offset);
        }

        public void ReadBytes(byte[] data, int offset)
        {

            using (MemoryStream memoryStream = new MemoryStream(data))
            {
                PacketReader packetReader = new PacketReader(memoryStream);
                memoryStream.Position = offset;
                this.WordCount = packetReader.ReadByte();
                this.AndXCommand = packetReader.ReadByte();
                this.AndXReserved = packetReader.ReadByte();
                this.AndXOffset = packetReader.ReadUInt16();
                this.MaxBufferSize = packetReader.ReadUInt16();
                this.MaxMpxCount = packetReader.ReadUInt16();
                this.VcNumber = packetReader.ReadUInt16();
                this.SessionKey = packetReader.ReadUInt32();
                this.SecurityBlobLength = packetReader.ReadUInt16();
                this.Reserved = packetReader.BigEndianReadUInt32();
                this.Capabilities = packetReader.ReadUInt32();
                this.ByteCount = packetReader.ReadUInt16();
                this.SecurityBlob = packetReader.ReadBytes(this.SecurityBlobLength);
            }

        }

    }
}
