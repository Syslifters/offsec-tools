using PingCastle.Data;

namespace PingCastleCommon.Data
{
    public class NetworkMapDataItem
    {
        public Subnet Network { get; set; }
        public string Source { get; set; }
        public string Name { get; set; }
        public string Description { get; set; }
        public string Location { get; set; }
        public string DomainFQDN { get; set; }
    }
}
