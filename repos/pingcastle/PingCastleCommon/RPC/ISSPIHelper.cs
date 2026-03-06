namespace PingCastleCommon.RPC;

using System.Net;

/// <summary>
/// Platform-agnostic interface for SSPI (Security Support Provider Interface) operations.
/// Isolates Windows-specific SSPI implementation details from cross-platform code.
/// </summary>
public interface ISSPIHelper
{
    string SecurityPackage { get; set; }

    void LoginClient(NetworkCredential login);
    void InitializeClient(out byte[] clientToken, byte[] serverToken, out bool bContinueProcessing);
    void InitializeServer(byte[] clientToken, out byte[] serverToken, out bool bContinueProcessing);
    void EncryptMessage(byte[] message, out byte[] encryptedBuffer);
    void DecryptMessage(int messageLength, byte[] encryptedBuffer, bool bUseClientContext, out byte[] decryptedBuffer);
    void SignMessage(byte[] message, out byte[] signedBuffer);
    void VerifyMessage(int messageLength, byte[] signedBuffer, out byte[] verifiedBuffer);
    byte[] GetSessionKey();
}
