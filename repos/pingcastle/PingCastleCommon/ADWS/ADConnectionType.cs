namespace PingCastle.ADWS;

public delegate void WorkOnReturnedObjectByADWS(ADItem Object);

public enum ADConnectionType
{
    Default = -1,
    ADWSThenLDAP = 0,
    ADWSOnly = 1,
    LDAPOnly = 2,
    LDAPThenADWS = 3,
    Unix = 4,
}