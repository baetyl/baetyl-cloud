INSERT INTO `baetyl_system_config` (`id`, `type`, `name`, `create_time`, `update_time`, `value`) VALUES
(2, 'address', 'node-address',  now(),  now(), 'https://host.docker.internal:30005'),
(3, 'address', 'active-address',  now(),  now(), 'https://0.0.0.0:30003');

INSERT INTO `baetyl_system_config` (`id`, `type`, `name`, `create_time`, `update_time`, `value`) VALUES
(14, 'baetyl-module', 'baetyl-init', now(), now(), 'hub.baidubce.com/baetyl/init:v2.0.0'),
(15, 'baetyl-module', 'baetyl-core', now(), now(), 'hub.baidubce.com/baetyl/core:v2.0.0'),
(16, 'baetyl-module', 'baetyl-function', now(), now(), 'hub.baidubce.com/baetyl/function:v2.0.0'),
(17, 'baetyl-module', 'baetyl-broker', now(), now(), 'hub.baidubce.com/baetyl/broker:v2.0.0'),
(29, 'resource', 'setup.sh', now(), now(), 'IyEvYmluL3NoCgpzZXQgLWUKCk9TPSQodW5hbWUpClRPS0VOPSJ7ey5Ub2tlbn19IgpSRVNPVVJDRV9VUkw9Int7LkJhZXR5bEhvc3R9fSIKU1VETz1zdWRvCgpleGVjX2NtZF9ub2JhaWwoKSB7CiAgZWNobyAiKyAkMiBiYXNoIC1jIFwiJDFcIiIKICAkMiBiYXNoIC1jICIkMSIKfQoKcHJpbnRfc3RhdHVzKCkgewogIGVjaG8gIiMjICQxIgp9Cgp1cmxfc2FmZV9jaGVjaygpIHsKICBpZiAhIGN1cmwgLUlmcyAkMSA+L2Rldi9udWxsOyB0aGVuCiAgICBwcmludF9zdGF0dXMgIkVSUk9SOiAkMSBpcyBpbnZhbGlkIG9yIFVucmVhY2hhYmxlISIKICBmaQp9CgpjaGVja19jbWQoKSB7CiAgY29tbWFuZCAtdiAkMSB8IGF3ayAne3ByaW50fScKfQoKZ2V0X2RlcGVuZGVuY2llcygpIHsKICBQUkVfSU5TVEFMTF9QS0dTPSIiCgogIGlmIFsgISAteCAiJChjaGVja19jbWQgY3VybCkiIF07IHRoZW4KICAgIFBSRV9JTlNUQUxMX1BLR1M9IiR7UFJFX0lOU1RBTExfUEtHU30gY3VybCIKICBmaQoKICBpZiBbICJYJHtQUkVfSU5TVEFMTF9QS0dTfSIgIT0gIlgiIF07IHRoZW4KICAgIGNhc2UgIiRPUyIgaW4KICAgIExpbnV4KQogICAgICBMU0JfRElTVD0kKC4gL2V0Yy9vcy1yZWxlYXNlICYmIGVjaG8gIiRJRCIgfCB0ciAnWzp1cHBlcjpdJyAnWzpsb3dlcjpdJykKICAgICAgY2FzZSAiJExTQl9ESVNUIiBpbgogICAgICB1YnVudHUgfCBkZWJpYW4gfCByYXNwYmlhbikKICAgICAgICBleGVjX2NtZF9ub2JhaWwgImFwdCB1cGRhdGUgJiYgYXB0IGluc3RhbGwgLS1uby1pbnN0YWxsLXJlY29tbWVuZHMgLXkgJHtQUkVfSU5TVEFMTF9QS0dTfSA+L2Rldi9udWxsIDI+JjEiICRTVURPCiAgICAgICAgOzsKICAgICAgY2VudG9zKQogICAgICAgIGV4ZWNfY21kX25vYmFpbCAieXVtIGluc3RhbGwgJHtQUkVfSU5TVEFMTF9QS0dTfSAteSA+L2Rldi9udWxsIDI+JjEiICRTVURPCiAgICAgICAgOzsKICAgICAgKikKICAgICAgICBwcmludF9zdGF0dXMgIllvdXIgT1MgaXMgbm90IHN1cHBvcnRlZCEiCiAgICAgICAgOzsKICAgICAgZXNhYwogICAgICA7OwogICAgRGFyd2luKQogICAgICBwcmludF9zdGF0dXMgIllvdSBtdXN0IGluc3RhbGwgJHtQUkVfSU5TVEFMTF9QS0dTfSB0byBjb250aW51ZS4uLiIKICAgICAgZXhpdCAwCiAgICAgIDs7CiAgICAqKQogICAgICBwcmludF9zdGF0dXMgIllvdXIgT1M6ICRPUyBpcyBub3Qgc3VwcG9ydGVkISIKICAgICAgZXhpdCAwCiAgICAgIDs7CiAgICBlc2FjCiAgZmkKfQoKaW5zdGFsbF9kb2NrZXIoKSB7CiAgVEFSR0VUX1VSTD1odHRwOi8vZ2V0LmRhb2Nsb3VkLmlvL2RvY2tlci8KICB1cmxfc2FmZV9jaGVjayAke1RBUkdFVF9VUkx9CiAgZXhlY19jbWRfbm9iYWlsICJjdXJsIC1zU0wgJHtUQVJHRVRfVVJMfSB8ICRTVURPIHNoIgoKICBpZiBbWyAhIC14ICIkKGNvbW1hbmQgLXYgZG9ja2VyKSIgXV07IHRoZW4KICAgICAgcHJpbnRfc3RhdHVzICJJbnN0YWxsIGRvY2tlciBmYWlsZWQhIENoZWNrIHRoZSBpbnN0YWxsaW5nIHByb2Nlc3MgZm9yIGhlbHAuLi4iCiAgZmkKCiAgaWYgW1sgISAteCAiJChjb21tYW5kIC12IHN5c3RlbWN0bCkiIF1dOyB0aGVuCiAgICAgIExTQl9ESVNUPSQoLiAvZXRjL29zLXJlbGVhc2UgJiYgZWNobyAiJElEIiB8IHRyICdbOnVwcGVyOl0nICdbOmxvd2VyOl0nKQogICAgICBjYXNlICIkTFNCX0RJU1QiIGluCiAgICAgIHVidW50dSB8IGRlYmlhbiB8IHJhc3BiaWFuKQogICAgICAgICAgZXhlY19jbWRfbm9iYWlsICJhcHQgdXBkYXRlICYmIGFwdCBpbnN0YWxsIC0tbm8taW5zdGFsbC1yZWNvbW1lbmRzIC15IHN5c3RlbWQgID4vZGV2L251bGwgMj4mMSIgJFNVRE8KICAgICAgICAgIDs7CiAgICAgIGNlbnRvcykKICAgICAgICAgIGV4ZWNfY21kX25vYmFpbCAieXVtIGluc3RhbGwgc3lzdGVtZCAteSAgPi9kZXYvbnVsbCAyPiYxIiAkU1VETwogICAgICAgICAgOzsKICAgICAgKikKICAgICAgICAgIHByaW50X3N0YXR1cyAiWW91ciBPUzogJE9TIGlzIG5vdCBzdXBwb3J0ZWQhIgogICAgICAgICAgZXhpdCAwCiAgICAgICAgICA7OwogICAgICBlc2FjCiAgZmkKCiAgZXhlY19jbWRfbm9iYWlsICJzeXN0ZW1jdGwgZW5hYmxlIGRvY2tlciIgJFNVRE8KICBleGVjX2NtZF9ub2JhaWwgInN5c3RlbWN0bCBzdGFydCBkb2NrZXIiICRTVURPCn0KCmNoZWNrX2FuZF9nZXRfa3ViZSgpIHsKICBpZiBbICEgLXggIiQoY2hlY2tfY21kIGt1YmVjdGwpIiBdOyB0aGVuCiAgICByZWFkIC1wICJLOFMvSzNTIGlzIG5vdCBpbnN0YWxsZWQgeWV0LCBkbyB5b3Ugd2FudCB1cyB0byBpbnN0YWxsIEszUyBmb3IgeW91PyBZZXMvTm8gKGRlZmF1bHQ6IFllcyk6IiBJU19JTlNUQUxMX0szUwogICAgaWYgWyAiJElTX0lOU1RBTExfSzNTIiA9ICJuIiAtbyAiJElTX0lOU1RBTExfSzNTIiA9ICJOIiAtbyAiJElTX0lOU1RBTExfSzNTIiA9ICJubyIgLW8gIiRJU19JTlNUQUxMX0szUyIgPSAiTk8iIF07IHRoZW4KICAgICAgZWNobyAiSzNTIGlzIG5lZWRlZCB0byBydW4gJHtOQU1FfSwgdGhpcyBzY3JpcHQgd2lsbCBleGl0IG5vdy4uLiIKICAgICAgZXhpdCAwCiAgICBmaQoKICAgIGlmIFsgJE9TID0gIkxpbnV4IiBdOyB0aGVuCiAgICAgICAgcmVhZCAtcCAiSzNTIGNvdWxkIHJ1biB3aXRoIGNvbnRhaW5lcmQvZG9ja2VyLCB3aGljaCBkbyB5b3Ugd2FudCB1cyB0byBpbnN0YWxsIGZvciB5b3U/IGNvbnRhaW5lcmQgZm9yIFllcywgZG9ja2VyIGZvciBObyAoZGVmYXVsdDogWWVzKToiIElTX0lOU1RBTExfQ09OVEFJTkVSRAogICAgICAgIGlmIFsgIiRJU19JTlNUQUxMX0NPTlRBSU5FUkQiID0gIm4iIC1vICIkSVNfSU5TVEFMTF9DT05UQUlORVJEIiA9ICJOIiAtbyAiJElTX0lOU1RBTExfQ09OVEFJTkVSRCIgPSAibm8iIC1vICIkSVNfSU5TVEFMTF9DT05UQUlORVJEIiA9ICJOTyIgXTsgdGhlbgogICAgICAgICAgaWYgWyAhIC14ICIkKGNoZWNrX2NtZCBkb2NrZXIpIiBdOyB0aGVuCiAgICAgICAgICAgIGluc3RhbGxfZG9ja2VyCiAgICAgICAgICBlbHNlCiAgICAgICAgICAgIHByaW50X3N0YXR1cyAiRG9ja2VyIGFscmVhZHkgaW5zdGFsbGVkIgogICAgICAgICAgZmkKICAgICAgICAgIGV4cG9ydCBJTlNUQUxMX0szU19FWEVDPSItLWRvY2tlciAtLXdyaXRlLWt1YmVjb25maWcgfi8ua3ViZS9jb25maWcgLS13cml0ZS1rdWJlY29uZmlnLW1vZGUgNjY2IgogICAgICAgIGVsc2UKICAgICAgICAgIGV4cG9ydCBJTlNUQUxMX0szU19FWEVDPSItLXdyaXRlLWt1YmVjb25maWcgfi8ua3ViZS9jb25maWcgLS13cml0ZS1rdWJlY29uZmlnLW1vZGUgNjY2IgogICAgICAgIGZpCgogICAgICBleGVjX2NtZF9ub2JhaWwgImN1cmwgLXNmTCBodHRwczovL2RvY3MucmFuY2hlci5jbi9rM3MvazNzLWluc3RhbGwuc2ggfCBJTlNUQUxMX0szU19NSVJST1I9Y24gc2ggLSIKCiAgICAgIGlmIFsgISAteCAiJChjaGVja19jbWQga3ViZWN0bCkiIF07IHRoZW4KICAgICAgICBwcmludF9zdGF0dXMgIkluc3RhbGwgazNzIGZhaWxlZCEgQ2hlY2sgdGhlIGluc3RhbGxpbmcgcHJvY2VzcyBmb3IgaGVscC4uLiIKICAgICAgICBleGl0IDAKICAgICAgZmkKCiAgICAgIGV4ZWNfY21kX25vYmFpbCAic3lzdGVtY3RsIGVuYWJsZSBrM3MiICRTVURPCiAgICAgIGV4ZWNfY21kX25vYmFpbCAic3lzdGVtY3RsIHN0YXJ0IGszcyIgJFNVRE8KCiAgICBlbGlmIFsgJE9TID0gIkRhcndpbiIgXTsgdGhlbgogICAgICBleGVjX2NtZF9ub2JhaWwgImN1cmwgLXNmTCBodHRwczovL2dldC5rM3MuaW8gfCBzaCAtIgogICAgZWxzZQogICAgICBwcmludF9zdGF0dXMgIldlIGFyZSBub3Qgc3VwcG9ydGluZyB5b3VyIHN5c3RlbSwgdGhpcyBzY3JpcHQgd2lsbCBleGl0IG5vdy4uLiIKICAgICAgZXhpdCAwCiAgICBmaQogIGZpCiAgZXhlY19jbWRfbm9iYWlsICJrdWJlY3RsIHZlcnNpb24iICRTVURPCn0KCmdldF9rdWJlX21hc3RlcigpIHsKICAkU1VETyBrdWJlY3RsIGdldCBubyB8IGF3ayAnL21hc3Rlci8gfHwgL2NvbnRyb2xwbGFuZS8nIHwgYXdrICd7cHJpbnQgJDF9Jwp9CgpjaGVja19iYWV0eWxfbmFtZXNwYWNlKCkgewogICRTVURPIGt1YmVjdGwgZ2V0IG5zIHwgZ3JlcCAnYmFldHlsLWVkZ2Utc3lzdGVtJyB8IGF3ayAne3ByaW50ICQxfScKfQoKY2hlY2tfYW5kX2luc3RhbGxfYmFldHlsKCkgewogIEJBRVRZTF9OQU1FU1BBQ0U9JChjaGVja19iYWV0eWxfbmFtZXNwYWNlKQogIGlmIFsgISAteiAiJEJBRVRZTF9OQU1FU1BBQ0UiIF07IHRoZW4KICAgIHJlYWQgLXAgIlRoZSBuYW1lc3BhY2UgJ2JhZXR5bC1lZGdlLXN5c3RlbScgYWxyZWFkeSBleGlzdHMsIGRvIHlvdSB3YW50IHRvIGNsZWFuIHVwIG9sZCBhcHBsaWNhdGlvbnMgYnkgZGVsZXRpbmcgdGhpcyBuYW1lc3BhY2U/IFllcy9ObyAoZGVmYXVsdDogWWVzKToiIElTX0RFTEVURV9OUwogICAgaWYgWyAiJElTX0RFTEVURV9OUyIgPSAibiIgLW8gIiRJU19ERUxFVEVfTlMiID0gIk4iIC1vICIkSVNfREVMRVRFX05TIiA9ICJubyIgLW8gIiRJU19ERUxFVEVfTlMiID0gIk5PIiBdOyB0aGVuCiAgICAgIGVjaG8gImJhZXR5bC1pbml0IGlzIG5vdCBpbnN0YWxsLCB0aGlzIHNjcmlwdCB3aWxsIGV4aXQgbm93Li4uIgogICAgICBleGl0IDAKICAgIGVsc2UKICAgICAgcmJhYz0kKCRTVURPIGt1YmVjdGwgZ2V0IGNsdXN0ZXJyb2xlYmluZGluZyB8IGdyZXAgYmFldHlsLWVkZ2Utc3lzdGVtLXJiYWMgfCBhd2sgJ3twcmludCAkMX0nKQogICAgICBpZiBbIC1uICIkcmJhYyIgXTsgdGhlbgogICAgICAgIGV4ZWNfY21kX25vYmFpbCAia3ViZWN0bCBkZWxldGUgY2x1c3RlcnJvbGViaW5kaW5nIGJhZXR5bC1lZGdlLXN5c3RlbS1yYmFjIiAkU1VETwogICAgICBmaQogICAgICBleGVjX2NtZF9ub2JhaWwgImt1YmVjdGwgZGVsZXRlIG5hbWVzcGFjZSBiYWV0eWwtZWRnZS1zeXN0ZW0iICRTVURPCiAgICBmaQogIGZpCgogIEtVQkVfTUFTVEVSX05PREVfTkFNRT0kKGdldF9rdWJlX21hc3RlcikKICBpZiBbICEgLXogIiRLVUJFX01BU1RFUl9OT0RFX05BTUUiIF07IHRoZW4KICAgIGV4ZWNfY21kX25vYmFpbCAibWtkaXIgLXAgLW0gNjY2IC92YXIvbGliL2JhZXR5bC9jb3JlLWRhdGEiICRTVURPCiAgICBleGVjX2NtZF9ub2JhaWwgIm1rZGlyIC1wIC1tIDY2NiAvdmFyL2xpYi9iYWV0eWwvYXBwLWRhdGEiICRTVURPCiAgICBleGVjX2NtZF9ub2JhaWwgIm1rZGlyIC1wIC1tIDY2NiAvdmFyL2xpYi9iYWV0eWwvY29yZS1zdG9yZSIgJFNVRE8KICAgIGV4ZWNfY21kX25vYmFpbCAibWtkaXIgLXAgLW0gNjY2IC92YXIvbG9nL2JhZXR5bC9jb3JlLWxvZyIgJFNVRE8KICAgIGV4ZWNfY21kX25vYmFpbCAibWtkaXIgLXAgLW0gNjY2IC92YXIvbGliL2JhZXR5bC9jb3JlLXBhZ2UiICRTVURPCiAgICBrdWJlX2FwcGx5ICIkUkVTT1VSQ0VfVVJML3YxL2FjdGl2ZS9iYWV0eWwtaW5pdC55bWw/dG9rZW49JFRPS0VOJm5vZGU9JEtVQkVfTUFTVEVSX05PREVfTkFNRSIKICBlbHNlCiAgICBwcmludF9zdGF0dXMgIkNhbiBub3QgZ2V0IGt1YmVybmV0ZXMgbWFzdGVyIG9yIGNvbnRyb2xwbGFuZSBub2RlLCB0aGlzIHNjcmlwdCB3aWxsIGV4aXQgbm93Li4uIgogIGZpCn0KCmt1YmVfYXBwbHkoKSB7CiAgVGVtcEZpbGU9JChta3RlbXAgdGVtcC5YWFhYWFgpCiAgZXhlY19jbWRfbm9iYWlsICJjdXJsIC1za2ZMIFwiJDFcIiA+JFRlbXBGaWxlIiAkU1VETwogIGV4ZWNfY21kX25vYmFpbCAia3ViZWN0bCBhcHBseSAtZiAkVGVtcEZpbGUiICRTVURPCiAgZXhlY19jbWRfbm9iYWlsICJybSAtZiAkVGVtcEZpbGUgMj4vZGV2L251bGwiICRTVURPCn0KCmNoZWNrX2FuZF9nZXRfbWV0cmljcygpIHsKICBNRVRSSUNTPSQoY2hlY2tfa3ViZV9yZXMgbWV0cmljcy1zZXJ2ZXIpCiAgaWYgWyAteiAiJE1FVFJJQ1MiIF07IHRoZW4KICAgIGt1YmVfYXBwbHkgIiRSRVNPVVJDRV9VUkwvdjEvYWN0aXZlL21ldHJpY3MueW1sIgogIGZpCn0KCmNoZWNrX2t1YmVfcmVzKCkgewogICRTVURPIGt1YmVjdGwgZ2V0IGRlcGxveW1lbnRzIC1BIHwgZ3JlcCAkMSB8IGF3ayAne3ByaW50ICQyfScKfQoKY2hlY2tfYW5kX2dldF9zdG9yYWdlKCkgewogIFBST1ZJU0lPTkVSPSQoY2hlY2tfa3ViZV9yZXMgbG9jYWwtcGF0aC1wcm92aXNpb25lcikKICBpZiBbIC16ICIkUFJPVklTSU9ORVIiIF07IHRoZW4KICAgIGt1YmVfYXBwbHkgIiRSRVNPVVJDRV9VUkwvdjEvYWN0aXZlL2xvY2FsLXBhdGgtc3RvcmFnZS55bWwiCiAgZmkKfQoKY2hlY2tfdXNlcigpIHsKICBpZiBbICQoaWQgLXUpIC1lcSAwIF07IHRoZW4KICAgIFNVRE89CiAgZmkKfQoKdW5pbnN0YWxsKCkgewogIGV4ZWNfY21kX25vYmFpbCAia3ViZWN0bCBkZWxldGUgbnMgYmFldHlsLWVkZ2UgYmFldHlsLWVkZ2Utc3lzdGVtIiAkU1VETwp9CgppbnN0YWxsKCkgewogIGNoZWNrX3VzZXIKICBnZXRfZGVwZW5kZW5jaWVzCiAgY2hlY2tfYW5kX2dldF9rdWJlCiAgY2hlY2tfYW5kX2dldF9tZXRyaWNzCiAgY2hlY2tfYW5kX2dldF9zdG9yYWdlCiAgY2hlY2tfYW5kX2luc3RhbGxfYmFldHlsCn0KCmNhc2UgQyIkMSIgaW4KQykKICBpbnN0YWxsCiAgOzsKQ3VuaW5zdGFsbCkKICB1bmluc3RhbGwKICA7OwpDaW5zdGFsbCkKICBpbnN0YWxsCiAgOzsKQyopCiAgVXNhZ2U6IHNldHVwLnNoIHsgaW5zdGFsbCB8IHVuaW5zdGFsbCB9CiAgOzsKZXNhYwoKZWNobyAiRG9uZSEiCmV4aXQgMAo='),
(30, 'resource', 'app-core.json', now(), now(), 'ewogICJuYW1lIjogInt7LkFwcE5hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInNlbGVjdG9yIjogImJhZXR5bC1ub2RlLW5hbWU9e3suTm9kZU5hbWV9fSIsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJ0eXBlIjogInt7LkFwcFR5cGV9fSIsCiAgInNlcnZpY2VzIjogWwogIHsKICAgICJuYW1lIjogImJhZXR5bC1jb3JlIiwKICAgICJpbWFnZSI6ICJ7ey5JbWFnZX19IiwKICAgICJyZXBsaWNhIjogMSwKICAgICJ2b2x1bWVNb3VudHMiOiBbCiAgICB7CiAgICAgICJuYW1lIjogImNvbmZpZyIsCiAgICAgICJtb3VudFBhdGgiOiAiL2V0Yy9iYWV0eWwiLAogICAgICAicmVhZE9ubHkiOiB0cnVlCiAgICB9LAogICAgewogICAgICAibmFtZSI6ICJkYXRhIiwKICAgICAgIm1vdW50UGF0aCI6ICIvdmFyL2xpYi9iYWV0eWwvY29yZS1kYXRhIgogICAgfSwKICAgIHsKICAgICAgICAgICJuYW1lIjogImFwcC1kYXRhIiwKICAgICAgICAgICJtb3VudFBhdGgiOiAiL3Zhci9saWIvYmFldHlsL2FwcC1kYXRhIgogICAgfSwKICAgIHsKICAgICAgIm5hbWUiOiAic3RvcmUiLAogICAgICAibW91bnRQYXRoIjogIi92YXIvbGliL2JhZXR5bC9zdG9yZSIKICAgIH0sCiAgICB7CiAgICAgICJuYW1lIjogImxvZyIsCiAgICAgICJtb3VudFBhdGgiOiAiL3Zhci9sb2cvYmFldHlsIgogICAgfSwKICAgIHsKICAgICAgIm5hbWUiOiAiY2VydC1zeW5jIiwKICAgICAgIm1vdW50UGF0aCI6ICIvdmFyL2xpYi9iYWV0eWwvY2VydC9zeW5jIgogICAgfQogICAgXSwKICAgICJwb3J0cyI6WwogICAgewogICAgICAiY29udGFpbmVyUG9ydCI6ODAsCiAgICAgICJob3N0UG9ydCI6IDMwMDUwLAogICAgICAicHJvdG9jb2wiOiJUQ1AiCiAgICB9CiAgICBdCiAgfQogIF0sCiAgInZvbHVtZXMiOiBbCiAgewogICAgIm5hbWUiOiAiY29uZmlnIiwKICAgICJjb25maWciOiB7CiAgICAgICJuYW1lIjogInt7LkNvbmZpZ05hbWV9fSIsCiAgICAgICJ2ZXJzaW9uIjogInt7LkNvbmZpZ1ZlcnNpb259fSIKICAgIH0KICB9LAogIHsKICAgICAgIm5hbWUiOiAiYXBwLWRhdGEiLAogICAgICAiaG9zdFBhdGgiOiB7CiAgICAgICAgInBhdGgiOiAiL3Zhci9saWIvYmFldHlsL2FwcC1kYXRhIgogICAgICB9CiAgfSwKICB7CiAgICAibmFtZSI6ICJkYXRhIiwKICAgICJob3N0UGF0aCI6IHsKICAgICAgInBhdGgiOiAiL3Zhci9saWIvYmFldHlsL2NvcmUtZGF0YSIKICAgIH0KICB9LAogIHsKICAgICAgIm5hbWUiOiAic3RvcmUiLAogICAgICAiaG9zdFBhdGgiOiB7CiAgICAgICAgInBhdGgiOiAiL3Zhci9saWIvYmFldHlsL2NvcmUtc3RvcmUiCiAgICAgIH0KICAgIH0sCiAgewogICAgIm5hbWUiOiAibG9nIiwKICAgICJob3N0UGF0aCI6IHsKICAgICAgInBhdGgiOiAiL3Zhci9saWIvYmFldHlsL2NvcmUtbG9nIgogICAgfQogIH0sCiAgewogICAgIm5hbWUiOiAiY2VydC1zeW5jIiwKICAgICJzZWNyZXQiOiB7CiAgICAgICJuYW1lIjogInt7LkNlcnRTeW5jfX0iLAogICAgICAidmVyc2lvbiI6ICJ7ey5DZXJ0U3luY1ZlcnNpb259fSIKICAgIH0KICB9CiAgXQp9'),
(31, 'resource', 'app-function.json', now(), now(), 'ewogICJuYW1lIjogInt7LkFwcE5hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInNlbGVjdG9yIjogImJhZXR5bC1ub2RlLW5hbWU9e3suTm9kZU5hbWV9fSIsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJ0eXBlIjogInt7LkFwcFR5cGV9fSIsCiAgInNlcnZpY2VzIjogWwogIHsKICAgICJuYW1lIjogImJhZXR5bC1mdW5jdGlvbiIsCiAgICAiaW1hZ2UiOiAie3suSW1hZ2V9fSIsCiAgICAicmVwbGljYSI6IDEsCiAgICAidm9sdW1lTW91bnRzIjogWwogICAgewogICAgICAibmFtZSI6ICJjb25maWciLAogICAgICAibW91bnRQYXRoIjogIi9ldGMvYmFldHlsIiwKICAgICAgInJlYWRPbmx5IjogdHJ1ZQogICAgfQogICAgXSwKICAgICJwb3J0cyI6WwogICAgewogICAgICAgICJjb250YWluZXJQb3J0Ijo4MCwKICAgICAgICAicHJvdG9jb2wiOiJUQ1AiCiAgICB9CiAgICBdCiAgfQogIF0sCiAgInZvbHVtZXMiOiBbCiAgewogICAgIm5hbWUiOiAiY29uZmlnIiwKICAgICJjb25maWciOiB7CiAgICAgICJuYW1lIjogInt7LkNvbmZpZ05hbWV9fSIsCiAgICAgICJ2ZXJzaW9uIjogInt7LkNvbmZpZ1ZlcnNpb259fSIKICAgIH0KICB9CiAgXQp9'),
(32, 'resource', 'baetyl-init.yml', now(), now(), 'LS0tCmFwaVZlcnNpb246IHYxCmtpbmQ6IE5hbWVzcGFjZQptZXRhZGF0YToKICBuYW1lOiB7ey5FZGdlU3lzdGVtTmFtZXNwYWNlfX0KCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBOYW1lc3BhY2UKbWV0YWRhdGE6CiAgbmFtZToge3suRWRnZU5hbWVzcGFjZX19CgotLS0KYXBpVmVyc2lvbjogdjEKa2luZDogU2VydmljZUFjY291bnQKbWV0YWRhdGE6CiAgbmFtZTogYmFldHlsLWVkZ2Utc3lzdGVtLXNlcnZpY2UtYWNjb3VudAogIG5hbWVzcGFjZToge3suRWRnZVN5c3RlbU5hbWVzcGFjZX19CgotLS0KIyBlbGV2YXRpb24gb2YgYXV0aG9yaXR5CmFwaVZlcnNpb246IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vdjFiZXRhMQpraW5kOiBDbHVzdGVyUm9sZUJpbmRpbmcKbWV0YWRhdGE6CiAgbmFtZTogYmFldHlsLWVkZ2Utc3lzdGVtLXJiYWMKc3ViamVjdHM6CiAgLSBraW5kOiBTZXJ2aWNlQWNjb3VudAogICAgbmFtZTogYmFldHlsLWVkZ2Utc3lzdGVtLXNlcnZpY2UtYWNjb3VudAogICAgbmFtZXNwYWNlOiB7ey5FZGdlU3lzdGVtTmFtZXNwYWNlfX0Kcm9sZVJlZjoKICBraW5kOiBDbHVzdGVyUm9sZQogIG5hbWU6IGNsdXN0ZXItYWRtaW4KICBhcGlHcm91cDogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pbwoKLS0tCmtpbmQ6IFN0b3JhZ2VDbGFzcwphcGlWZXJzaW9uOiBzdG9yYWdlLms4cy5pby92MQptZXRhZGF0YToKICBuYW1lOiBsb2NhbC1zdG9yYWdlCnByb3Zpc2lvbmVyOiBrdWJlcm5ldGVzLmlvL25vLXByb3Zpc2lvbmVyCnZvbHVtZUJpbmRpbmdNb2RlOiBXYWl0Rm9yRmlyc3RDb25zdW1lcgoKe3stIGlmIC5DZXJ0U3luY1BlbX19Ci0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZWNyZXQKbWV0YWRhdGE6CiAgbmFtZToge3suQ2VydFN5bmN9fQogIG5hbWVzcGFjZToge3suRWRnZVN5c3RlbU5hbWVzcGFjZX19CnR5cGU6IE9wYXF1ZQpkYXRhOgogIGNsaWVudC5wZW06ICd7ey5DZXJ0U3luY1BlbX19JwogIGNsaWVudC5rZXk6ICd7ey5DZXJ0U3luY0tleX19JwogIGNhLnBlbTogJ3t7LkNlcnRTeW5jQ2F9fScKe3stIGVuZH19Cgp7ey0gaWYgLkNlcnRBY3RpdmVDYX19Ci0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZWNyZXQKbWV0YWRhdGE6CiAgbmFtZToge3suQ2VydEFjdGl2ZX19CiAgbmFtZXNwYWNlOiB7ey5FZGdlU3lzdGVtTmFtZXNwYWNlfX0KdHlwZTogT3BhcXVlCmRhdGE6CiAgY2EucGVtOiAne3suQ2VydEFjdGl2ZUNhfX0nCnt7LSBlbmR9fQoKLS0tCiMgYmFldHlsLWluaXQgY29uZmlnbWFwCmFwaVZlcnNpb246IHYxCmtpbmQ6IENvbmZpZ01hcAptZXRhZGF0YToKICBuYW1lOiBiYWV0eWwtaW5pdC1jb25maWcKICBuYW1lc3BhY2U6IHt7LkVkZ2VTeXN0ZW1OYW1lc3BhY2V9fQpkYXRhOgogIHNlcnZpY2UueW1sOiB8LQogICAgZW5naW5lOgogICAgICBrdWJlcm5ldGVzOgogICAgICAgIGluQ2x1c3RlcjogdHJ1ZQogICAgc3luYzoKICAgICAgY2xvdWQ6CiAgICAgICAgaHR0cDoKICAgICAgICAgIGFkZHJlc3M6IHt7Lk5vZGVBZGRyZXNzfX0KICAgICAgICAgIGNhOiB2YXIvbGliL2JhZXR5bC9jZXJ0L2NhLnBlbQogICAgICAgICAga2V5OiB2YXIvbGliL2JhZXR5bC9jZXJ0L2NsaWVudC5rZXkKICAgICAgICAgIGNlcnQ6IHZhci9saWIvYmFldHlsL2NlcnQvY2xpZW50LnBlbQogICAgICAgICAgaW5zZWN1cmVTa2lwVmVyaWZ5OiB0cnVlCiAgICAgIGVkZ2U6CiAgICAgICAgZG93bmxvYWRQYXRoOiAvdmFyL2xpYi9iYWV0eWwvY29yZS1kYXRhCiAgICB7ey0gaWYgLkNlcnRTeW5jUGVtfX0KICAgIHt7ZWxzZX19CiAgICBpbml0OgogICAgICBiYXRjaDoKICAgICAgICBuYW1lOiB7ey5CYXRjaE5hbWV9fQogICAgICAgIG5hbWVzcGFjZToge3suTmFtZXNwYWNlfX0KICAgICAgICBzZWN1cml0eVR5cGU6IHt7LlNlY3VyaXR5VHlwZX19CiAgICAgICAgc2VjdXJpdHlLZXk6IHt7LlNlY3VyaXR5S2V5fX0KICAgICAgY2xvdWQ6CiAgICAgICAgaHR0cDoKICAgICAgICAgIGFkZHJlc3M6IHt7LkFjdGl2ZUFkZHJlc3N9fQogICAgICAgICAgY2E6IHZhci9saWIvYmFldHlsL2NlcnQtYWN0aXZlL2NhLnBlbQogICAgICAgICAgaW5zZWN1cmVTa2lwVmVyaWZ5OiB0cnVlCiAgICAgIGVkZ2U6CiAgICAgICAgZG93bmxvYWRQYXRoOiAvdmFyL2xpYi9iYWV0eWwvY29yZS1kYXRhCiAgICAgIGFjdGl2ZToKICAgICAgICB7ey0gaWYgZXEgLlByb29mVHlwZSAiaW5wdXQifX0KICAgICAgICBzZXJ2ZXI6CiAgICAgICAgICBsaXN0ZW46IDAuMC4wLjA6e3suSG9zdFBvcnR9fQogICAgICAgICAgcGFnZXM6IC92YXIvbGliL2JhZXR5bC9wYWdlCiAgICAgICAge3stIGVuZH19CiAgICAgICAgZmluZ2VycHJpbnRzOgogICAgICAgICAgLSBwcm9vZjoge3suUHJvb2ZUeXBlfX0KICAgICAgICAgICAgdmFsdWU6IHt7LlByb29mVmFsdWV9fQogICAgICAgIGF0dHJpYnV0ZXM6CiAgICAgICAgICAtIG5hbWU6IGJhdGNoCiAgICAgICAgICAgIGxhYmVsOiBCYXRjaE5hbWUKICAgICAgICAgICAgdmFsdWU6IHt7LkJhdGNoTmFtZX19CiAgICAgICAgICAtIG5hbWU6IG5hbWVzcGFjZQogICAgICAgICAgICBsYWJlbDogTmFtZXNwYWNlCiAgICAgICAgICAgIHZhbHVlOiB7ey5OYW1lc3BhY2V9fQogICAgICAgICAge3stIGlmIGVxIC5Qcm9vZlR5cGUgImlucHV0In19CiAgICAgICAgICAtIG5hbWU6IHNuCiAgICAgICAgICAgIGxhYmVsOiBTZXJpYWxOdW1iZXIKICAgICAgICAgIC0gbmFtZTogZGVzY3JpcHRpb24KICAgICAgICAgICAgbGFiZWw6IERlc2NyaXB0aW9uCiAgICAgICAgICB7ey0gZW5kfX0KICAgIHt7LSBlbmR9fQogICAgbG9nZ2VyOgogICAgICBmaWxlbmFtZTogdmFyL2xvZy9iYWV0eWwvc2VydmljZS5sb2cKICAgICAgbGV2ZWw6IGRlYnVnCgotLS0KIyBiYWV0eWwtaW5pdCBoZWFkbGVzcyBzZXJ2aWNlCmFwaVZlcnNpb246IHYxCmtpbmQ6IFNlcnZpY2UKbWV0YWRhdGE6CiAgbmFtZTogYmFldHlsLWluaXQKICBuYW1lc3BhY2U6IHt7LkVkZ2VTeXN0ZW1OYW1lc3BhY2V9fQogIGxhYmVsczoKICAgIGJhZXR5bC1zZXJ2aWNlLW5hbWU6IGJhZXR5bC1pbml0CiAgICBiYWV0eWwtYXBwLW5hbWU6IGJhZXR5bC1pbml0CnNwZWM6CiAgcHVibGlzaE5vdFJlYWR5QWRkcmVzc2VzOiB0cnVlCiAge3stIGlmIC5Qcm9vZlR5cGV9fQogIHt7LSBpZiBlcSAuUHJvb2ZUeXBlICJpbnB1dCJ9fQogIHR5cGU6IE5vZGVQb3J0CiAgcG9ydHM6CiAgICAtIHBvcnQ6IHt7LkNvbnRhaW5lclBvcnR9fQogICAgICB0YXJnZXRQb3J0OiB7ey5Db250YWluZXJQb3J0fX0KICAgICAgbm9kZVBvcnQ6IHt7Lkhvc3RQb3J0fX0KICB7e2Vsc2V9fQogIGNsdXN0ZXJJUDogTm9uZQogIHt7LSBlbmR9fQogIHt7ZWxzZX19CiAgY2x1c3RlcklQOiBOb25lCiAge3stIGVuZH19CiAgc2VsZWN0b3I6CiAgICBiYWV0eWwtc2VydmljZS1uYW1lOiBiYWV0eWwtaW5pdAogICAgCi0tLQojIGJhZXR5bC1pbml0IGRlcGxveW1lbnQKYXBpVmVyc2lvbjogYXBwcy92MQpraW5kOiBEZXBsb3ltZW50Cm1ldGFkYXRhOgogIG5hbWU6IGJhZXR5bC1pbml0CiAgbmFtZXNwYWNlOiB7ey5FZGdlU3lzdGVtTmFtZXNwYWNlfX0KICBsYWJlbHM6CiAgICBiYWV0eWwtYXBwLW5hbWU6IGJhZXR5bC1pbml0CiAgICBiYWV0eWwtc2VydmljZS1uYW1lOiBiYWV0eWwtaW5pdApzcGVjOgogIHNlbGVjdG9yOgogICAgbWF0Y2hMYWJlbHM6CiAgICAgIGJhZXR5bC1zZXJ2aWNlLW5hbWU6IGJhZXR5bC1pbml0CiAgcmVwbGljYXM6IDEKICB0ZW1wbGF0ZToKICAgIG1ldGFkYXRhOgogICAgICBsYWJlbHM6CiAgICAgICAgYmFldHlsLWFwcC1uYW1lOiBiYWV0eWwtaW5pdAogICAgICAgIGJhZXR5bC1zZXJ2aWNlLW5hbWU6IGJhZXR5bC1pbml0CiAgICBzcGVjOgogICAgICBub2RlTmFtZToge3suS3ViZU5vZGVOYW1lfX0KICAgICAgc2VydmljZUFjY291bnROYW1lOiBiYWV0eWwtZWRnZS1zeXN0ZW0tc2VydmljZS1hY2NvdW50CiAgICAgIGNvbnRhaW5lcnM6CiAgICAgICAgLSBuYW1lOiBiYWV0eWwtaW5pdAogICAgICAgICAgaW1hZ2U6IHt7LkltYWdlfX0KICAgICAgICAgIGltYWdlUHVsbFBvbGljeTogSWZOb3RQcmVzZW50CiAgICAgICAgICB7ey0gaWYgLlByb29mVHlwZX19CiAgICAgICAgICB7ey0gaWYgZXEgLlByb29mVHlwZSAiaW5wdXQifX0KICAgICAgICAgIHBvcnRzOgogICAgICAgICAgICAtIGNvbnRhaW5lclBvcnQ6IHt7LkNvbnRhaW5lclBvcnR9fQogICAgICAgICAge3stIGVuZH19CiAgICAgICAgICB7ey0gZW5kfX0KICAgICAgICAgIGVudjoKICAgICAgICAgICAgLSBuYW1lOiBLVUJFX05PREVfTkFNRQogICAgICAgICAgICAgIHZhbHVlRnJvbToKICAgICAgICAgICAgICAgIGZpZWxkUmVmOgogICAgICAgICAgICAgICAgICBmaWVsZFBhdGg6IHNwZWMubm9kZU5hbWUKICAgICAgICAgIHZvbHVtZU1vdW50czoKICAgICAgICAgICAgLSBuYW1lOiBkYXRhCiAgICAgICAgICAgICAgbW91bnRQYXRoOiAvdmFyL2xpYi9iYWV0eWwvY29yZS1kYXRhCiAgICAgICAgICAgIC0gbmFtZTogc3RvcmUKICAgICAgICAgICAgICBtb3VudFBhdGg6IC92YXIvbGliL2JhZXR5bC9zdG9yZQogICAgICAgICAgICAtIG5hbWU6IGxvZwogICAgICAgICAgICAgIG1vdW50UGF0aDogdmFyL2xvZy9iYWV0eWwKICAgICAgICAgICAge3stIGlmIC5DZXJ0U3luY1BlbX19CiAgICAgICAgICAgIC0gbmFtZTogY2VydC1zeW5jCiAgICAgICAgICAgICAgbW91bnRQYXRoOiB2YXIvbGliL2JhZXR5bC9jZXJ0CiAgICAgICAgICAgIHt7LSBlbmR9fQogICAgICAgICAgICB7ey0gaWYgLkNlcnRBY3RpdmVDYX19CiAgICAgICAgICAgIC0gbmFtZTogY2VydC1hY3RpdmUKICAgICAgICAgICAgICBtb3VudFBhdGg6IHZhci9saWIvYmFldHlsL2NlcnQtYWN0aXZlCiAgICAgICAgICAgIHt7LSBlbmR9fQogICAgICAgICAgICB7ey0gaWYgLlByb29mVHlwZX19CiAgICAgICAgICAgIHt7LSBpZiBlcSAuUHJvb2ZUeXBlICJzbiJ9fQogICAgICAgICAgICAtIG5hbWU6IHNuCiAgICAgICAgICAgICAgbW91bnRQYXRoOiAvdmFyL2xpYi9iYWV0eWwvc24KICAgICAgICAgICAge3stIGVuZH19CiAgICAgICAgICAgIHt7LSBlbmR9fQogICAgICAgICAgICAtIG5hbWU6IGNvbmZpZwogICAgICAgICAgICAgIG1vdW50UGF0aDogL2V0Yy9iYWV0eWwKICAgICAgdm9sdW1lczoKICAgICAgICAtIG5hbWU6IGRhdGEKICAgICAgICAgIGhvc3RQYXRoOgogICAgICAgICAgICBwYXRoOiAvdmFyL2xpYi9iYWV0eWwvY29yZS1kYXRhCiAgICAgICAgLSBuYW1lOiBzdG9yZQogICAgICAgICAgaG9zdFBhdGg6CiAgICAgICAgICAgIHBhdGg6IC92YXIvbGliL2JhZXR5bC9jb3JlLXN0b3JlCiAgICAgICAgLSBuYW1lOiBsb2cKICAgICAgICAgIGhvc3RQYXRoOgogICAgICAgICAgICBwYXRoOiAvdmFyL2xvZy9iYWV0eWwvY29yZS1sb2cKICAgICAgICB7ey0gaWYgLkNlcnRTeW5jUGVtfX0KICAgICAgICAtIG5hbWU6IGNlcnQtc3luYwogICAgICAgICAgc2VjcmV0OgogICAgICAgICAgICBzZWNyZXROYW1lOiB7ey5DZXJ0U3luY319CiAgICAgICAge3stIGVuZH19CiAgICAgICAge3stIGlmIC5DZXJ0QWN0aXZlQ2F9fQogICAgICAgIC0gbmFtZTogY2VydC1hY3RpdmUKICAgICAgICAgIHNlY3JldDoKICAgICAgICAgICAgc2VjcmV0TmFtZToge3suQ2VydEFjdGl2ZX19CiAgICAgICAge3stIGVuZH19CiAgICAgICAge3stIGlmIC5Qcm9vZlR5cGV9fQogICAgICAgIHt7LSBpZiBlcSAuUHJvb2ZUeXBlICJzbiJ9fQogICAgICAgIC0gbmFtZTogc24KICAgICAgICAgIGhvc3RQYXRoOgogICAgICAgICAgICBwYXRoOiB7ey5Tbkhvc3RQYXRofX0KICAgICAgICB7ey0gZW5kfX0KICAgICAgICB7ey0gZW5kfX0KICAgICAgICAtIG5hbWU6IGNvbmZpZwogICAgICAgICAgY29uZmlnTWFwOgogICAgICAgICAgICBuYW1lOiBiYWV0eWwtaW5pdC1jb25maWc='),
(33, 'resource', 'config-core.json', now(), now(), 'ewogICJuYW1lIjogInt7LkNvbmZpZ05hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInN5c3RlbSI6IHRydWUsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtYXBwLW5hbWUiOiJ7ey5BcHBOYW1lfX0iLAogICAgImJhZXR5bC1ub2RlLW5hbWUiOiJ7ey5Ob2RlTmFtZX19IiwKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJkYXRhIjogewogICAgInNlcnZpY2UueW1sIjogImVuZ2luZTpcbiAga3ViZXJuZXRlczpcbiAgICBpbkNsdXN0ZXI6IHRydWVcbnN5bmM6XG4gIGNsb3VkOlxuICAgIGh0dHA6XG4gICAgICBhZGRyZXNzOiB7ey5Ob2RlQWRkcmVzc319XG4gICAgICBjYTogL3Zhci9saWIvYmFldHlsL2NlcnQvc3luYy9jYS5wZW1cbiAgICAgIGtleTogL3Zhci9saWIvYmFldHlsL2NlcnQvc3luYy9jbGllbnQua2V5XG4gICAgICBjZXJ0OiAvdmFyL2xpYi9iYWV0eWwvY2VydC9zeW5jL2NsaWVudC5wZW1cbiAgICAgIGluc2VjdXJlU2tpcFZlcmlmeTogdHJ1ZVxuICBlZGdlOlxuICAgIGRvd25sb2FkUGF0aDogL3Zhci9saWIvYmFldHlsL2NvcmUtZGF0YVxubG9nZ2VyOlxuICBmaWxlbmFtZTogdmFyL2xvZy9iYWV0eWwvc2VydmljZS5sb2dcbiAgbGV2ZWw6IGRlYnVnIgogIH0KfQ=='),
(34, 'resource', 'config-function.json', now(), now(), 'ewogICJuYW1lIjogInt7LkNvbmZpZ05hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInN5c3RlbSI6IHRydWUsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtYXBwLW5hbWUiOiJ7ey5BcHBOYW1lfX0iLAogICAgImJhZXR5bC1ub2RlLW5hbWUiOiJ7ey5Ob2RlTmFtZX19IiwKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJkYXRhIjogewogICAgInNlcnZpY2UueW1sIjogIiIKICB9Cn0='),
(35, 'resource', 'metrics.yml', now(), now(), 'a2luZDogQ2x1c3RlclJvbGUKYXBpVmVyc2lvbjogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby92MQptZXRhZGF0YToKICBuYW1lOiBzeXN0ZW06YWdncmVnYXRlZC1tZXRyaWNzLXJlYWRlcgogIGxhYmVsczoKICAgIHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vYWdncmVnYXRlLXRvLXZpZXc6ICJ0cnVlIgogICAgcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby9hZ2dyZWdhdGUtdG8tZWRpdDogInRydWUiCiAgICByYmFjLmF1dGhvcml6YXRpb24uazhzLmlvL2FnZ3JlZ2F0ZS10by1hZG1pbjogInRydWUiCnJ1bGVzOgogIC0gYXBpR3JvdXBzOiBbIm1ldHJpY3MuazhzLmlvIl0KICAgIHJlc291cmNlczogWyJwb2RzIiwgIm5vZGVzIl0KICAgIHZlcmJzOiBbImdldCIsICJsaXN0IiwgIndhdGNoIl0KLS0tCmFwaVZlcnNpb246IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vdjFiZXRhMQpraW5kOiBDbHVzdGVyUm9sZUJpbmRpbmcKbWV0YWRhdGE6CiAgbmFtZTogbWV0cmljcy1zZXJ2ZXI6c3lzdGVtOmF1dGgtZGVsZWdhdG9yCnJvbGVSZWY6CiAgYXBpR3JvdXA6IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8KICBraW5kOiBDbHVzdGVyUm9sZQogIG5hbWU6IHN5c3RlbTphdXRoLWRlbGVnYXRvcgpzdWJqZWN0czoKICAtIGtpbmQ6IFNlcnZpY2VBY2NvdW50CiAgICBuYW1lOiBtZXRyaWNzLXNlcnZlcgogICAgbmFtZXNwYWNlOiBrdWJlLXN5c3RlbQotLS0KYXBpVmVyc2lvbjogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby92MWJldGExCmtpbmQ6IFJvbGVCaW5kaW5nCm1ldGFkYXRhOgogIG5hbWU6IG1ldHJpY3Mtc2VydmVyLWF1dGgtcmVhZGVyCiAgbmFtZXNwYWNlOiBrdWJlLXN5c3RlbQpyb2xlUmVmOgogIGFwaUdyb3VwOiByYmFjLmF1dGhvcml6YXRpb24uazhzLmlvCiAga2luZDogUm9sZQogIG5hbWU6IGV4dGVuc2lvbi1hcGlzZXJ2ZXItYXV0aGVudGljYXRpb24tcmVhZGVyCnN1YmplY3RzOgogIC0ga2luZDogU2VydmljZUFjY291bnQKICAgIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCi0tLQphcGlWZXJzaW9uOiBhcGlyZWdpc3RyYXRpb24uazhzLmlvL3YxYmV0YTEKa2luZDogQVBJU2VydmljZQptZXRhZGF0YToKICBuYW1lOiB2MWJldGExLm1ldHJpY3MuazhzLmlvCnNwZWM6CiAgc2VydmljZToKICAgIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCiAgZ3JvdXA6IG1ldHJpY3MuazhzLmlvCiAgdmVyc2lvbjogdjFiZXRhMQogIGluc2VjdXJlU2tpcFRMU1ZlcmlmeTogdHJ1ZQogIGdyb3VwUHJpb3JpdHlNaW5pbXVtOiAxMDAKICB2ZXJzaW9uUHJpb3JpdHk6IDEwMAotLS0KYXBpVmVyc2lvbjogdjEKa2luZDogU2VydmljZUFjY291bnQKbWV0YWRhdGE6CiAgbmFtZTogbWV0cmljcy1zZXJ2ZXIKICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCi0tLQphcGlWZXJzaW9uOiBhcHBzL3YxCmtpbmQ6IERlcGxveW1lbnQKbWV0YWRhdGE6CiAgbmFtZTogbWV0cmljcy1zZXJ2ZXIKICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCiAgbGFiZWxzOgogICAgazhzLWFwcDogbWV0cmljcy1zZXJ2ZXIKc3BlYzoKICBzZWxlY3RvcjoKICAgIG1hdGNoTGFiZWxzOgogICAgICBrOHMtYXBwOiBtZXRyaWNzLXNlcnZlcgogIHRlbXBsYXRlOgogICAgbWV0YWRhdGE6CiAgICAgIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgICAgIGxhYmVsczoKICAgICAgICBrOHMtYXBwOiBtZXRyaWNzLXNlcnZlcgogICAgc3BlYzoKICAgICAgc2VydmljZUFjY291bnROYW1lOiBtZXRyaWNzLXNlcnZlcgogICAgICB2b2x1bWVzOgogICAgICAgICMgbW91bnQgaW4gdG1wIHNvIHdlIGNhbiBzYWZlbHkgdXNlIGZyb20tc2NyYXRjaCBpbWFnZXMgYW5kL29yIHJlYWQtb25seSBjb250YWluZXJzCiAgICAgICAgLSBuYW1lOiB0bXAtZGlyCiAgICAgICAgICBlbXB0eURpcjoge30KICAgICAgY29udGFpbmVyczoKICAgICAgICAtIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgICAgICAgICBpbWFnZTogJ3JhbmNoZXIvbWV0cmljcy1zZXJ2ZXI6djAuMy42JwogICAgICAgICAgaW1hZ2VQdWxsUG9saWN5OiBJZk5vdFByZXNlbnQKICAgICAgICAgIGNvbW1hbmQ6CiAgICAgICAgICAgIC0gL21ldHJpY3Mtc2VydmVyCiAgICAgICAgICAgIC0gLS1rdWJlbGV0LWluc2VjdXJlLXRscwogICAgICAgICAgICAtIC0ta3ViZWxldC1wcmVmZXJyZWQtYWRkcmVzcy10eXBlcz1JbnRlcm5hbEROUyxJbnRlcm5hbElQLEV4dGVybmFsRE5TLEV4dGVybmFsSVAsSG9zdG5hbWUKICAgICAgICAgIHZvbHVtZU1vdW50czoKICAgICAgICAgICAgLSBuYW1lOiB0bXAtZGlyCiAgICAgICAgICAgICAgbW91bnRQYXRoOiAvdG1wCi0tLQphcGlWZXJzaW9uOiB2MQpraW5kOiBTZXJ2aWNlCm1ldGFkYXRhOgogIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgbmFtZXNwYWNlOiBrdWJlLXN5c3RlbQogIGxhYmVsczoKICAgIGt1YmVybmV0ZXMuaW8vbmFtZTogIk1ldHJpY3Mtc2VydmVyIgogICAga3ViZXJuZXRlcy5pby9jbHVzdGVyLXNlcnZpY2U6ICJ0cnVlIgpzcGVjOgogIHNlbGVjdG9yOgogICAgazhzLWFwcDogbWV0cmljcy1zZXJ2ZXIKICBwb3J0czoKICAgIC0gcG9ydDogNDQzCiAgICAgIHByb3RvY29sOiBUQ1AKICAgICAgdGFyZ2V0UG9ydDogNDQzCi0tLQphcGlWZXJzaW9uOiByYmFjLmF1dGhvcml6YXRpb24uazhzLmlvL3YxCmtpbmQ6IENsdXN0ZXJSb2xlCm1ldGFkYXRhOgogIG5hbWU6IHN5c3RlbTptZXRyaWNzLXNlcnZlcgpydWxlczoKICAtIGFwaUdyb3VwczoKICAgICAgLSAiIgogICAgcmVzb3VyY2VzOgogICAgICAtIHBvZHMKICAgICAgLSBub2RlcwogICAgICAtIG5vZGVzL3N0YXRzCiAgICAgIC0gbmFtZXNwYWNlcwogICAgdmVyYnM6CiAgICAgIC0gZ2V0CiAgICAgIC0gbGlzdAogICAgICAtIHdhdGNoCi0tLQphcGlWZXJzaW9uOiByYmFjLmF1dGhvcml6YXRpb24uazhzLmlvL3YxCmtpbmQ6IENsdXN0ZXJSb2xlQmluZGluZwptZXRhZGF0YToKICBuYW1lOiBzeXN0ZW06bWV0cmljcy1zZXJ2ZXIKcm9sZVJlZjoKICBhcGlHcm91cDogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pbwogIGtpbmQ6IENsdXN0ZXJSb2xlCiAgbmFtZTogc3lzdGVtOm1ldHJpY3Mtc2VydmVyCnN1YmplY3RzOgogIC0ga2luZDogU2VydmljZUFjY291bnQKICAgIG5hbWU6IG1ldHJpY3Mtc2VydmVyCiAgICBuYW1lc3BhY2U6IGt1YmUtc3lzdGVtCg=='),
(36, 'resource', 'local-path-storage.yml', now(), now(), 'YXBpVmVyc2lvbjogdjEKa2luZDogTmFtZXNwYWNlCm1ldGFkYXRhOgogIG5hbWU6IGxvY2FsLXBhdGgtc3RvcmFnZQotLS0KYXBpVmVyc2lvbjogdjEKa2luZDogU2VydmljZUFjY291bnQKbWV0YWRhdGE6CiAgbmFtZTogbG9jYWwtcGF0aC1wcm92aXNpb25lci1zZXJ2aWNlLWFjY291bnQKICBuYW1lc3BhY2U6IGxvY2FsLXBhdGgtc3RvcmFnZQotLS0KYXBpVmVyc2lvbjogcmJhYy5hdXRob3JpemF0aW9uLms4cy5pby92MQpraW5kOiBDbHVzdGVyUm9sZQptZXRhZGF0YToKICBuYW1lOiBsb2NhbC1wYXRoLXByb3Zpc2lvbmVyLXJvbGUKcnVsZXM6CiAgLSBhcGlHcm91cHM6IFsiIl0KICAgIHJlc291cmNlczogWyJub2RlcyIsICJwZXJzaXN0ZW50dm9sdW1lY2xhaW1zIl0KICAgIHZlcmJzOiBbImdldCIsICJsaXN0IiwgIndhdGNoIl0KICAtIGFwaUdyb3VwczogWyIiXQogICAgcmVzb3VyY2VzOiBbImVuZHBvaW50cyIsICJwZXJzaXN0ZW50dm9sdW1lcyIsICJwb2RzIl0KICAgIHZlcmJzOiBbIioiXQogIC0gYXBpR3JvdXBzOiBbIiJdCiAgICByZXNvdXJjZXM6IFsiZXZlbnRzIl0KICAgIHZlcmJzOiBbImNyZWF0ZSIsICJwYXRjaCJdCiAgLSBhcGlHcm91cHM6IFsic3RvcmFnZS5rOHMuaW8iXQogICAgcmVzb3VyY2VzOiBbInN0b3JhZ2VjbGFzc2VzIl0KICAgIHZlcmJzOiBbImdldCIsICJsaXN0IiwgIndhdGNoIl0KLS0tCmFwaVZlcnNpb246IHJiYWMuYXV0aG9yaXphdGlvbi5rOHMuaW8vdjEKa2luZDogQ2x1c3RlclJvbGVCaW5kaW5nCm1ldGFkYXRhOgogIG5hbWU6IGxvY2FsLXBhdGgtcHJvdmlzaW9uZXItYmluZApyb2xlUmVmOgogIGFwaUdyb3VwOiByYmFjLmF1dGhvcml6YXRpb24uazhzLmlvCiAga2luZDogQ2x1c3RlclJvbGUKICBuYW1lOiBsb2NhbC1wYXRoLXByb3Zpc2lvbmVyLXJvbGUKc3ViamVjdHM6CiAgLSBraW5kOiBTZXJ2aWNlQWNjb3VudAogICAgbmFtZTogbG9jYWwtcGF0aC1wcm92aXNpb25lci1zZXJ2aWNlLWFjY291bnQKICAgIG5hbWVzcGFjZTogbG9jYWwtcGF0aC1zdG9yYWdlCi0tLQphcGlWZXJzaW9uOiBhcHBzL3YxCmtpbmQ6IERlcGxveW1lbnQKbWV0YWRhdGE6CiAgbmFtZTogbG9jYWwtcGF0aC1wcm92aXNpb25lcgogIG5hbWVzcGFjZTogbG9jYWwtcGF0aC1zdG9yYWdlCnNwZWM6CiAgcmVwbGljYXM6IDEKICBzZWxlY3RvcjoKICAgIG1hdGNoTGFiZWxzOgogICAgICBhcHA6IGxvY2FsLXBhdGgtcHJvdmlzaW9uZXIKICB0ZW1wbGF0ZToKICAgIG1ldGFkYXRhOgogICAgICBsYWJlbHM6CiAgICAgICAgYXBwOiBsb2NhbC1wYXRoLXByb3Zpc2lvbmVyCiAgICBzcGVjOgogICAgICBzZXJ2aWNlQWNjb3VudE5hbWU6IGxvY2FsLXBhdGgtcHJvdmlzaW9uZXItc2VydmljZS1hY2NvdW50CiAgICAgIGNvbnRhaW5lcnM6CiAgICAgICAgLSBuYW1lOiBsb2NhbC1wYXRoLXByb3Zpc2lvbmVyCiAgICAgICAgICBpbWFnZTogcmFuY2hlci9sb2NhbC1wYXRoLXByb3Zpc2lvbmVyOnYwLjAuMTIKICAgICAgICAgIGltYWdlUHVsbFBvbGljeTogSWZOb3RQcmVzZW50CiAgICAgICAgICBjb21tYW5kOgogICAgICAgICAgICAtIGxvY2FsLXBhdGgtcHJvdmlzaW9uZXIKICAgICAgICAgICAgLSAtLWRlYnVnCiAgICAgICAgICAgIC0gc3RhcnQKICAgICAgICAgICAgLSAtLWNvbmZpZwogICAgICAgICAgICAtIC9ldGMvY29uZmlnL2NvbmZpZy5qc29uCiAgICAgICAgICB2b2x1bWVNb3VudHM6CiAgICAgICAgICAgIC0gbmFtZTogY29uZmlnLXZvbHVtZQogICAgICAgICAgICAgIG1vdW50UGF0aDogL2V0Yy9jb25maWcvCiAgICAgICAgICBlbnY6CiAgICAgICAgICAgIC0gbmFtZTogUE9EX05BTUVTUEFDRQogICAgICAgICAgICAgIHZhbHVlRnJvbToKICAgICAgICAgICAgICAgIGZpZWxkUmVmOgogICAgICAgICAgICAgICAgICBmaWVsZFBhdGg6IG1ldGFkYXRhLm5hbWVzcGFjZQogICAgICB2b2x1bWVzOgogICAgICAgIC0gbmFtZTogY29uZmlnLXZvbHVtZQogICAgICAgICAgY29uZmlnTWFwOgogICAgICAgICAgICBuYW1lOiBsb2NhbC1wYXRoLWNvbmZpZwotLS0KYXBpVmVyc2lvbjogc3RvcmFnZS5rOHMuaW8vdjEKa2luZDogU3RvcmFnZUNsYXNzCm1ldGFkYXRhOgogIG5hbWU6IGxvY2FsLXBhdGgKcHJvdmlzaW9uZXI6IHJhbmNoZXIuaW8vbG9jYWwtcGF0aAp2b2x1bWVCaW5kaW5nTW9kZTogV2FpdEZvckZpcnN0Q29uc3VtZXIKcmVjbGFpbVBvbGljeTogRGVsZXRlCi0tLQpraW5kOiBDb25maWdNYXAKYXBpVmVyc2lvbjogdjEKbWV0YWRhdGE6CiAgbmFtZTogbG9jYWwtcGF0aC1jb25maWcKICBuYW1lc3BhY2U6IGxvY2FsLXBhdGgtc3RvcmFnZQpkYXRhOgogIGNvbmZpZy5qc29uOiB8LQogICAgewogICAgICAgICAgICAibm9kZVBhdGhNYXAiOlsKICAgICAgICAgICAgewogICAgICAgICAgICAgICAgICAgICJub2RlIjoiREVGQVVMVF9QQVRIX0ZPUl9OT05fTElTVEVEX05PREVTIiwKICAgICAgICAgICAgICAgICAgICAicGF0aHMiOlsiL29wdC9sb2NhbC1wYXRoLXByb3Zpc2lvbmVyIl0KICAgICAgICAgICAgfQogICAgICAgICAgICBdCiAgICB9'),
(37, 'resource', 'app-broker.json', now(), now(), 'ewogICJuYW1lIjogInt7LkFwcE5hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInNlbGVjdG9yIjogImJhZXR5bC1ub2RlLW5hbWU9e3suTm9kZU5hbWV9fSIsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJ0eXBlIjogInt7LkFwcFR5cGV9fSIsCiAgInNlcnZpY2VzIjogWwogIHsKICAgICJuYW1lIjogImJhZXR5bC1icm9rZXIiLAogICAgImltYWdlIjogInt7LkltYWdlfX0iLAogICAgInJlcGxpY2EiOiAxLAogICAgInZvbHVtZU1vdW50cyI6IFsKICAgIHsKICAgICAgIm5hbWUiOiAiY29uZmlnIiwKICAgICAgIm1vdW50UGF0aCI6ICIvZXRjL2JhZXR5bCIsCiAgICAgICJyZWFkT25seSI6IHRydWUKICAgIH0KICAgIF0sCiAgICAicG9ydHMiOlsKICAgIHsKICAgICAgICAiY29udGFpbmVyUG9ydCI6ODAsCiAgICAgICAgInByb3RvY29sIjoiVENQIgogICAgfQogICAgXQogIH0KICBdLAogICJ2b2x1bWVzIjogWwogIHsKICAgICJuYW1lIjogImNvbmZpZyIsCiAgICAiY29uZmlnIjogewogICAgICAibmFtZSI6ICJ7ey5Db25maWdOYW1lfX0iLAogICAgICAidmVyc2lvbiI6ICJ7ey5Db25maWdWZXJzaW9ufX0iCiAgICB9CiAgfQogIF0KfQ=='),
(38, 'resource', 'config-broker.json', now(), now(), 'ewogICJuYW1lIjogInt7LkNvbmZpZ05hbWV9fSIsCiAgIm5hbWVzcGFjZSI6ICJ7ey5OYW1lc3BhY2V9fSIsCiAgInN5c3RlbSI6IHRydWUsCiAgImxhYmVscyI6IHsKICAgICJiYWV0eWwtYXBwLW5hbWUiOiJ7ey5BcHBOYW1lfX0iLAogICAgImJhZXR5bC1ub2RlLW5hbWUiOiJ7ey5Ob2RlTmFtZX19IiwKICAgICJiYWV0eWwtY2xvdWQtc3lzdGVtIjoidHJ1ZSIKICB9LAogICJkYXRhIjogewogICAgInNlcnZpY2UueW1sIjogIiIKICB9Cn0='),
(40, 'baetyl_version', 'latest',  now(), now(), 'v2.0.0'),
(50, 'certificate', 'baetyl.ca',  now(), now(), '8a1064640312470dae855829fe8d74dd'),
(60, 'object', 'object-source',  now(), now(), 'baidubos');

INSERT INTO `baetyl_certificate` (`cert_id`, `parent_id`, `type`, `common_name`, `description`, `csr`, `content`, `not_before`, `not_after`, `private_key`) VALUES
('8a1064640312470dae855829fe8d74dd', '', 'TypeIssuingCA', 'root.ca', 'cloud/ca.pem', 'LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJQ1lqQ0NBZ2lnQXdJQkFnSURBWWFpTUFvR0NDcUdTTTQ5QkFNQ01JR2xNUXN3Q1FZRFZRUUdFd0pEVGpFUQpNQTRHQTFVRUNCTUhRbVZwYW1sdVp6RVpNQmNHQTFVRUJ4TVFTR0ZwWkdsaGJpQkVhWE4wY21samRERVZNQk1HCkExVUVDUk1NUW1GcFpIVWdRMkZ0Y0hWek1ROHdEUVlEVlFRUkV3WXhNREF3T1RNeEhqQWNCZ05WQkFvVEZVeHAKYm5WNElFWnZkVzVrWVhScGIyNGdSV1JuWlRFUE1BMEdBMVVFQ3hNR1FrRkZWRmxNTVJBd0RnWURWUVFERXdkeQpiMjkwTG1OaE1DQVhEVEl3TURNeU56QTVOVGMxTkZvWUR6SXdOVEF3TXpJM01EazFOelUwV2pDQnBURUxNQWtHCkExVUVCaE1DUTA0eEVEQU9CZ05WQkFnVEIwSmxhV3BwYm1jeEdUQVhCZ05WQkFjVEVFaGhhV1JwWVc0Z1JHbHoKZEhKcFkzUXhGVEFUQmdOVkJBa1RERUpoYVdSMUlFTmhiWEIxY3pFUE1BMEdBMVVFRVJNR01UQXdNRGt6TVI0dwpIQVlEVlFRS0V4Vk1hVzUxZUNCR2IzVnVaR0YwYVc5dUlFVmtaMlV4RHpBTkJnTlZCQXNUQmtKQlJWUlpUREVRCk1BNEdBMVVFQXhNSGNtOXZkQzVqWVRCWk1CTUdCeXFHU000OUFnRUdDQ3FHU000OUF3RUhBMElBQk91WUhKWTAKODNBZVhXQVI0NWxiUlJ6SUNkSnpveWwwek9GSUs2V1c3WHVKV0QvSGhZeENlZFlQYXVtZ3dWSTdSUnNOYis2MQpWcFM5NWFGTUNNRFc0eXFqSXpBaE1BNEdBMVVkRHdFQi93UUVBd0lCaGpBUEJnTlZIUk1CQWY4RUJUQURBUUgvCk1Bb0dDQ3FHU000OUJBTUNBMGdBTUVVQ0lHaEJzaXMvUklhNCttekJlRHdySUNBY1laaGpxUVJjOHBQMTcyZ1EKeVZSdkFpRUF5UnpqN2tsbi94WGQxeVBXNkJFL2NhdkNwNnRtd3B2QURVbGFYUzVGRTlNPQotLS0tLUVORCBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0K', 'LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlDWWpDQ0FnaWdBd0lCQWdJREFZYWlNQW9HQ0NxR1NNNDlCQU1DTUlHbE1Rc3dDUVlEVlFRR0V3SkRUakVRDQpNQTRHQTFVRUNCTUhRbVZwYW1sdVp6RVpNQmNHQTFVRUJ4TVFTR0ZwWkdsaGJpQkVhWE4wY21samRERVZNQk1HDQpBMVVFQ1JNTVFtRnBaSFVnUTJGdGNIVnpNUTh3RFFZRFZRUVJFd1l4TURBd09UTXhIakFjQmdOVkJBb1RGVXhwDQpiblY0SUVadmRXNWtZWFJwYjI0Z1JXUm5aVEVQTUEwR0ExVUVDeE1HUWtGRlZGbE1NUkF3RGdZRFZRUURFd2R5DQpiMjkwTG1OaE1DQVhEVEl3TURNeU56QTVOVGMxTkZvWUR6SXdOVEF3TXpJM01EazFOelUwV2pDQnBURUxNQWtHDQpBMVVFQmhNQ1EwNHhFREFPQmdOVkJBZ1RCMEpsYVdwcGJtY3hHVEFYQmdOVkJBY1RFRWhoYVdScFlXNGdSR2x6DQpkSEpwWTNReEZUQVRCZ05WQkFrVERFSmhhV1IxSUVOaGJYQjFjekVQTUEwR0ExVUVFUk1HTVRBd01Ea3pNUjR3DQpIQVlEVlFRS0V4Vk1hVzUxZUNCR2IzVnVaR0YwYVc5dUlFVmtaMlV4RHpBTkJnTlZCQXNUQmtKQlJWUlpUREVRDQpNQTRHQTFVRUF4TUhjbTl2ZEM1allUQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJPdVlISlkwDQo4M0FlWFdBUjQ1bGJSUnpJQ2RKem95bDB6T0ZJSzZXVzdYdUpXRC9IaFl4Q2VkWVBhdW1nd1ZJN1JSc05iKzYxDQpWcFM5NWFGTUNNRFc0eXFqSXpBaE1BNEdBMVVkRHdFQi93UUVBd0lCaGpBUEJnTlZIUk1CQWY4RUJUQURBUUgvDQpNQW9HQ0NxR1NNNDlCQU1DQTBnQU1FVUNJR2hCc2lzL1JJYTQrbXpCZUR3cklDQWNZWmhqcVFSYzhwUDE3MmdRDQp5VlJ2QWlFQXlSemo3a2xuL3hYZDF5UFc2QkUvY2F2Q3A2dG13cHZBRFVsYVhTNUZFOU09DQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tDQo=', '2020-03-27 09:57:54', '2050-03-27 09:57:54', 'LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tDQpNSGNDQVFFRUlJMHFaQ1FNSjNRRDZGbWhENVBFeEU1cW5kUU5oMXRWNWZReTdKR3BIYkNIb0FvR0NDcUdTTTQ5DQpBd0VIb1VRRFFnQUU2NWdjbGpUemNCNWRZQkhqbVZ0RkhNZ0owbk9qS1hUTTRVZ3JwWmJ0ZTRsWVA4ZUZqRUo1DQoxZzlxNmFEQlVqdEZHdzF2N3JWV2xMM2xvVXdJd05iaktnPT0NCi0tLS0tRU5EIEVDIFBSSVZBVEUgS0VZLS0tLS0NCg==');
