CISCO-HSRP-EXT-MIB::cHsrpExtIfTrackedEntry:
  mib: CISCO-HSRP-EXT-MIB
  object: cHsrpExtIfTrackedEntry
  index:
    - type: 
      oid: .
      name:  # 
      syntax: 
  augments: 
  discovery_attribute: 
  attributes:
    cHsrpExtIfTrackedPriority:
      oid: .1.3.6.1.4.1.9.9.107.1.1.1.1.2
      name: cHsrpExtIfTrackedPriority
      syntax: Unsigned32
    cHsrpExtIfTrackedIpNone:
      oid: .1.3.6.1.4.1.9.9.107.1.1.1.1.4
      name: cHsrpExtIfTrackedIpNone
      syntax: TruthValue

CISCO-HSRP-EXT-MIB::cHsrpExtIfStandbyEntry:
  mib: CISCO-HSRP-EXT-MIB
  object: cHsrpExtIfStandbyEntry
  index:
    - type: 
      oid: .
      name:  # 
      syntax: 
  augments: 
  discovery_attribute: 
  attributes:
    cHsrpExtIfStandbyDestAddr:
      oid: .1.3.6.1.4.1.9.9.107.1.1.3.1.3
      name: cHsrpExtIfStandbyDestAddr
      syntax: IpAddressNoSuffix
    cHsrpExtIfStandbySourceAddr:
      oid: .1.3.6.1.4.1.9.9.107.1.1.3.1.5
      name: cHsrpExtIfStandbySourceAddr
      syntax: IpAddressNoSuffix

CISCO-HSRP-EXT-MIB::cHsrpExtIfEntry:
  mib: CISCO-HSRP-EXT-MIB
  object: cHsrpExtIfEntry
  index:
    - type: 
      oid: .
      name:  # 
      syntax: 
  augments: 
  discovery_attribute: 
  attributes:
    cHsrpExtIfUseBIA:
      oid: .1.3.6.1.4.1.9.9.107.1.2.1.1.1
      name: cHsrpExtIfUseBIA
      syntax: TruthValue
