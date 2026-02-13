using System;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Security;
using System.Net.Sockets;
using System.Security.Authentication;
using System.Security.Authentication.ExtendedProtection;
using System.Security.Cryptography.X509Certificates;

namespace PingCastle.misc
{
    public enum ConnectionTesterStatus
    {
        InvalidData,
        InitializationFailed,
        NoAuthenticationNeeded,
        AuthenticationFailure,
        InternalError,
        ChannelBindingEnabled,
        ChannelBindingDisabled,
        SignatureRequired,
        SignatureNotRequired,
        AuthenticationSuccessfull,
        LocalCall
    }

    abstract class ConnectionTester : IDisposable
    {

        public static ConnectionTesterStatus TestExtendedAuthentication(Uri uri, NetworkCredential credential, string logPrefix)
        {
            try
            {
                using (var tester = CreateTester(uri))
                {
                    tester.LogPrefix = logPrefix;
                    tester.Credential = credential;
                    return tester.TestExtendedAuthentication(uri);
                }
            }
            catch (LocalCallException)
            {
                return ConnectionTesterStatus.LocalCall;
            }
            catch (Exception ex)
            {
                Trace.WriteLine(logPrefix + "Exception: " + ex.Message);
                Trace.WriteLine(logPrefix + "StackTrace: " + ex.StackTrace);
                return ConnectionTesterStatus.InvalidData;
            }
        }

        public static ConnectionTesterStatus TestSignatureRequiredEnabled(Uri uri, NetworkCredential credential, string logPrefix)
        {
            try
            {
                using (var tester = CreateTester(uri))
                {
                    tester.LogPrefix = logPrefix;
                    tester.Credential = credential;
                    return tester.TestSignatureRequiredEnabled(uri);
                }
            }
            catch (LocalCallException)
            {
                return ConnectionTesterStatus.LocalCall;
            }
            catch (Exception ex)
            {
                Trace.WriteLine(logPrefix + "Exception: " + ex.Message);
                Trace.WriteLine(logPrefix + "StackTrace: " + ex.StackTrace);
                return ConnectionTesterStatus.InvalidData;
            }
        }

        public static ConnectionTesterStatus TestConnection(Uri uri, NetworkCredential credential, string logPrefix)
        {
            try
            {
                using (var tester = CreateTester(uri))
                {
                    tester.LogPrefix = logPrefix;
                    tester.Credential = credential;
                    return tester.Test(uri, false, true);
                }
            }
            catch (LocalCallException)
            {
                return ConnectionTesterStatus.LocalCall;
            }
            catch (Exception ex)
            {
                Trace.WriteLine(logPrefix + "Exception: " + ex.Message);
                Trace.WriteLine(logPrefix + "StackTrace: " + ex.StackTrace);
                return ConnectionTesterStatus.InvalidData;
            }
        }

        private class LocalCallException : Exception
        {

        }

        static bool AcceptEveryServerCertificate(object sender, X509Certificate certificate, X509Chain chain, SslPolicyErrors sslPolicyErrors)
        {
            return true;
        }


        static ConnectionTester CreateTester(Uri uri)
        {
            if (uri.Scheme == "ldaps" || uri.Scheme == "ldap")
            {
                return new ConnectionTesterLdap();
            }
            if (uri.Scheme == "https" || uri.Scheme == "http")
            {
                return new ConnectionTesterHttp();
            }
            Trace.WriteLine("CreateTester: Invalid scheme " + uri.Scheme);
            throw new ArgumentOutOfRangeException("uri", uri, "Scheme not supported");
        }

        protected ConnectionTesterStatus TestExtendedAuthentication(Uri uri)
        {
            Trace.WriteLine(LogPrefix + "Testing " + uri);

            Trace.WriteLine(LogPrefix + "testing WITH Extended Protection: ");
            var WithExtended = Test(uri, false, true);
            Trace.WriteLine(WithExtended.ToString());

            if (WithExtended != ConnectionTesterStatus.AuthenticationSuccessfull)
            {
                return WithExtended;
            }

            Trace.WriteLine(LogPrefix + "testing WITHOUT Extended Protection: ");
            var WithoutExtended = Test(uri, false, false);
            Trace.WriteLine(WithoutExtended.ToString());

            if (WithoutExtended != ConnectionTesterStatus.AuthenticationSuccessfull && WithoutExtended != ConnectionTesterStatus.AuthenticationFailure)
            {
                return WithExtended;
            }
            if (WithoutExtended == ConnectionTesterStatus.AuthenticationSuccessfull)
            {
                return ConnectionTesterStatus.ChannelBindingDisabled;
            }
            return ConnectionTesterStatus.ChannelBindingEnabled;
        }

        protected ConnectionTesterStatus TestSignatureRequiredEnabled(Uri uri)
        {
            Trace.WriteLine(LogPrefix + "Testing " + uri);

            Trace.WriteLine(LogPrefix + "testing WITH signature required: ");
            var WithSignature = Test(uri, false);
            Trace.WriteLine(WithSignature.ToString());

            if (WithSignature != ConnectionTesterStatus.AuthenticationSuccessfull)
            {
                return WithSignature;
            }

            Trace.WriteLine(LogPrefix + "testing WITHOUT signature required: ");
            var WithoutSignature = Test(uri, true);
            Trace.WriteLine(WithoutSignature.ToString());

            if (WithoutSignature != ConnectionTesterStatus.AuthenticationSuccessfull && WithoutSignature != ConnectionTesterStatus.AuthenticationFailure)
            {
                return WithSignature;
            }
            if (WithoutSignature == ConnectionTesterStatus.AuthenticationSuccessfull)
            {
                return ConnectionTesterStatus.SignatureNotRequired;
            }
            return ConnectionTesterStatus.SignatureRequired;
        }

        ConnectionTesterStatus Test(Uri uri, bool disableSigning = false, bool enableChannelBinding = false)
        {
            try
            {
                int port = uri.Port;
                if (port == -1)
                {
                    port = (uri.Scheme == "ldaps" ? 636 : 389);
                }
                // default LDAPS port is 636
                using (TcpClient tcpclient = new TcpClient(uri.Host, port))
                {
                    tcpclient.SendTimeout = 1000;
                    tcpclient.ReceiveTimeout = 1000;

                    if (uri.Scheme.EndsWith("s"))
                    {
                        using (SslStream sslStream = new SslStream(tcpclient.GetStream(), false, AcceptEveryServerCertificate, null))
                        {
                            // force all auth protocol.
                            // if SslProtocols.None used (default is not specified), if Tls12 is added in Default and if the only supported protocol on server side, exception will occure:
                            // The client and server cannot communicate, because they do not possess a common algorithm
                            sslStream.AuthenticateAsClient(uri.Host, null, SslProtocols.Ssl2 | SslProtocols.Ssl3 | SslProtocols.Tls | SslProtocols.Tls11 | SslProtocols.Tls12, false);

                            if (enableChannelBinding)
                            {
                                var tc = sslStream.TransportContext;
                                var cb = tc.GetChannelBinding(ChannelBindingKind.Endpoint);
                                InitializeNTAuthentication(true, cb);
                            }
                            else
                            {
                                InitializeNTAuthentication(true);
                            }

                            return SendPackets(sslStream, uri);
                        }
                    }
                    else
                    {
                        InitializeNTAuthentication(disableSigning);
                        return SendPackets(tcpclient.GetStream(), uri);
                    }
                }
            }
            catch (Exception ex)
            {
                Trace.WriteLine(LogPrefix + "ExceptionType: " + ex.GetType());
                Trace.WriteLine(LogPrefix + "Exception: " + ex.Message);
                Trace.WriteLine(LogPrefix + "StackTrace: " + ex.StackTrace);
                return ConnectionTesterStatus.InternalError;
            }

        }

        protected abstract ConnectionTesterStatus SendPackets(Stream stream, Uri uri);

        static int GetNTLMSSPOffset(byte[] output)
        {
            for (int i = 0; i < output.Length - 24; i++)
            {
                if (output[i] == 'N' && output[i + 1] == 'T' && output[i + 2] == 'L' && output[i + 3] == 'M' &&
                    output[i + 4] == 'S' && output[i + 5] == 'S' && output[i + 6] == 'P' && output[i + 7] == '\0')
                {
                    return i;
                }
            }
            return -1;
        }

        protected string package = "Negotiate";
        private NegotiateAuthentication _negotiateAuth;
        private bool _currentDisableSigning;
        private ChannelBinding _currentChannelBinding;

        public string LogPrefix { get; private set; }

        private NetworkCredential Credential;

        protected byte[] GetOutgoingBlob(byte[] incomingBlob)
        {
            ReadOnlySpan<byte> incomingSpan = incomingBlob ?? new byte[0];
            var status = _negotiateAuth.GetOutgoingBlob(incomingSpan, out NegotiateAuthenticationStatusCode statusCode);
            byte[] response = status?.ToArray() ?? new byte[0];

            int offset = GetNTLMSSPOffset(response);
            if (offset > 0)
            {
                var messageType = BitConverter.ToInt32(response, offset + 8);
                if (messageType == 2)
                {
                    var flag = BitConverter.ToInt32(response, offset + 20);
                    if ((flag & 0x00004000) != 0)
                    {
                        Trace.WriteLine(LogPrefix + "Local CALL");
                        throw new LocalCallException();
                    }
                }
            }
            return response;
        }

        protected void InitializeNTAuthentication(bool disableSigning = false, ChannelBinding channelBinding = null)
        {
            _currentDisableSigning = disableSigning;
            _currentChannelBinding = channelBinding;

            _negotiateAuth?.Dispose();

            var options = new NegotiateAuthenticationClientOptions
            {
                Package = package,
                Binding = channelBinding,
                RequiredProtectionLevel = disableSigning ? ProtectionLevel.None : ProtectionLevel.Sign
            };

            if (Credential != null)
            {
                options.Credential = Credential;
            }

            _negotiateAuth = new NegotiateAuthentication(options);
        }

        protected void ReinitializeNTAuthentication()
        {
            _negotiateAuth?.Dispose();

            var options = new NegotiateAuthenticationClientOptions
            {
                Package = package,
                Binding = _currentChannelBinding,
                RequiredProtectionLevel = _currentDisableSigning ? ProtectionLevel.None : ProtectionLevel.Sign
            };

            if (Credential != null)
            {
                options.Credential = Credential;
            }

            _negotiateAuth = new NegotiateAuthentication(options);
        }

        public void Dispose()
        {
            _negotiateAuth?.Dispose();
        }
    }
}
