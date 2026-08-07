using PingCastle.Data;

namespace PingCastleCommon.Data
{
    public class NetworkMapDataView
    {
        public int Order { get; set; }
        public Subnet FrameNetwork { get; set; }
        public bool HasData { get; set; }
        public int RecordCount { get; set; }
        public int ForestCount { get; set; }
    }
}
