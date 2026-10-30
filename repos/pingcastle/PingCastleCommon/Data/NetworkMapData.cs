using System;
using System.Collections.Generic;
using PingCastle.ADWS;

namespace PingCastleCommon.Data
{
    public class NetworkMapData
    {
        public List<NetworkMapDataView> Views { get; set; } = new();
        public Dictionary<string, List<NetworkMapDataItem>> NetworkRange { get; set; } = new();
        public List<NetworkMapDCItem> DomainControllers { get; set; } = new();
    }
}
