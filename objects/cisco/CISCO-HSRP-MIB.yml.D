CISCO-HSRP-MIB::cHsrpGlobalConfig:
  mib: CISCO-HSRP-MIB
  object: cHsrpGlobalConfig
  index:
    - type: 
      oid: .
      name:  # 
      syntax: 
  augments: 
  discovery_attribute: 
  attributes:
    cHsrpConfigTimeout:
      oid: .1.3.6.1.4.1.9.9.106.1.1.1
      name: cHsrpConfigTimeout
      syntax: Unsigned32

CISCO-HSRP-MIB::cHsrpGrpEntry:
  mib: CISCO-HSRP-MIB
  object: cHsrpGrpEntry
  index:
    - type: 
      oid: .
      name:  # 
      syntax: 
  augments: 
  discovery_attribute: 
  attributes:
    cHsrpGrpAuth:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.2
      name: cHsrpGrpAuth
      syntax: DisplayString
    cHsrpGrpPriority:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.3
      name: cHsrpGrpPriority
      syntax: Unsigned32
    cHsrpGrpPreempt:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.4
      name: cHsrpGrpPreempt
      syntax: TruthValue
    cHsrpGrpPreemptDelay:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.5
      name: cHsrpGrpPreemptDelay
      syntax: Unsigned32
    cHsrpGrpUseConfiguredTimers:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.6
      name: cHsrpGrpUseConfiguredTimers
      syntax: TruthValue
    cHsrpGrpConfiguredHelloTime:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.7
      name: cHsrpGrpConfiguredHelloTime
      syntax: Unsigned32
    cHsrpGrpConfiguredHoldTime:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.8
      name: cHsrpGrpConfiguredHoldTime
      syntax: Unsigned32
    cHsrpGrpLearnedHelloTime:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.9
      name: cHsrpGrpLearnedHelloTime
      syntax: Unsigned32
    cHsrpGrpLearnedHoldTime:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.10
      name: cHsrpGrpLearnedHoldTime
      syntax: Unsigned32
    cHsrpGrpVirtualIpAddr:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.11
      name: cHsrpGrpVirtualIpAddr
      syntax: IpAddress
    cHsrpGrpUseConfigVirtualIpAddr:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.12
      name: cHsrpGrpUseConfigVirtualIpAddr
      syntax: TruthValue
    cHsrpGrpActiveRouter:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.13
      name: cHsrpGrpActiveRouter
      syntax: IpAddress
    cHsrpGrpStandbyRouter:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.14
      name: cHsrpGrpStandbyRouter
      syntax: IpAddress
    cHsrpGrpStandbyState:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.15
      name: cHsrpGrpStandbyState
      syntax: HsrpState
    cHsrpGrpVirtualMacAddr:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.16
      name: cHsrpGrpVirtualMacAddr
      syntax: MacAddress
    cHsrpGrpIpNone:
      oid: .1.3.6.1.4.1.9.9.106.1.2.1.1.18
      name: cHsrpGrpIpNone
      syntax: TruthValue
