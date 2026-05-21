namespace PingCastleCommon.RPC;

using System.Net;

/// <summary>
/// No-op implementation of ISSPIHelper for non-Windows platforms.
/// Provides graceful degradation when SSPI operations are unavailable.
/// </summary>
public class NullSSPIHelper : ISSPIHelper
{
    public string SecurityPackage { get; set; }

    public void LoginClient(NetworkCredential login)
    {
    }

    public void InitializeClient(out byte[] clientToken, byte[] serverToken, out bool bContinueProcessing)
    {
        clientToken = null;
        bContinueProcessing = false;
    }

    public void InitializeServer(byte[] clientToken, out byte[] serverToken, out bool bContinueProcessing)
    {
        serverToken = null;
        bContinueProcessing = false;
    }

    public void EncryptMessage(byte[] message, out byte[] encryptedBuffer)
    {
        encryptedBuffer = null;
    }

    public void DecryptMessage(int messageLength, byte[] encryptedBuffer, bool bUseClientContext, out byte[] decryptedBuffer)
    {
        decryptedBuffer = null;
    }

    public void SignMessage(byte[] message, out byte[] signedBuffer)
    {
        signedBuffer = null;
    }

    public void VerifyMessage(int messageLength, byte[] signedBuffer, out byte[] verifiedBuffer)
    {
        verifiedBuffer = null;
    }

    public byte[] GetSessionKey()
    {
        return null;
    }
}
