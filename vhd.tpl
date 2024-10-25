{{define "tb" -}}
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;

ENTITY tb_{{ .Entity }} IS
END tb_{{ .Entity }};

ARCHITECTURE behavior OF tb_{{ .Entity }} IS
  -- Component Declaration for the Unit Under Test (UUT)
  COMPONENT {{ .Entity }} IS
  PORT(
    {{ range $i, $port := .Ports -}}
    {{ if eq $i (sub (len $.Ports) 1) -}}
    {{ if eq $port.Type "std_logic" -}}
    {{ $port.Name }} : {{ $port.InOut }} {{ $port.Type }}
    {{ else -}}
    {{ $port.Name }} : {{ $port.InOut }} {{ $port.Type }}({{ $port.MSB }} downto {{ $port.LSB }})
    {{ end -}}
    {{ else -}}
    {{ if eq $port.Type "std_logic" -}}
    {{ $port.Name }} : {{ $port.InOut }} {{ $port.Type }};
    {{ else -}}
    {{ $port.Name }} : {{ $port.InOut }} {{ $port.Type }}({{ $port.MSB }} downto {{ $port.LSB }});
    {{ end -}}
    {{ end -}}
    {{- end -}}
  );
  END COMPONENT;
  --Inputs
  {{ range $i, $port := .Ports -}}
  {{ if eq $port.InOut "in" -}}
  {{ if eq $port.Type "std_logic" -}}
  signal {{ $port.Name }} : {{ $port.Type }} := '0';
  {{ else -}}
  signal {{ $port.Name }} : {{ $port.Type }}({{ $port.MSB }} downto {{ $port.LSB }}) := (others => '0');
  {{ end -}}
  {{ end }}
  {{- end }}
  --Outputs
  {{ range $i, $port := .Ports -}}
  {{ if eq $port.InOut "out" -}}
  {{ if eq $port.Type "std_logic" -}}
  signal {{ $port.Name }} : {{ $port.Type }};
  {{ else -}}
  signal {{ $port.Name }} : {{ $port.Type }}({{ $port.MSB }} downto {{ $port.LSB }});
  {{ end -}}
  {{ end }}
  {{- end }}
  -- Clock period definitions
  constant CLK50M_period : time := 20 ns;
BEGIN
  -- Instantiate the Unit Under Test (UUT)
  uut: {{ .Entity }} PORT MAP (
    {{ range $i, $port := .Ports -}}
    {{ if eq $i (sub (len $.Ports) 1) -}}
    {{ $port.Name }} => {{ $port.Name }}
    {{ else -}}
    {{ $port.Name }} => {{ $port.Name }},
    {{ end -}}
    {{ end -}}
  );

  -- Clock process definitions
  CLK50M_process :process
  begin
		{{ .ClkPort.Name }} <= '0';
		wait for CLK50M_period/2;
		{{ .ClkPort.Name }} <= '1';
		wait for CLK50M_period/2;
  end process;

  -- Stimulus process
  stim_proc: process
  begin
    -- hold reset state for 100 ns.
	  {{ .ResetPort.Name }} <= '1';
    wait for 100 ns;
    {{ .ResetPort.Name }} <= '0';
    
    -- insert stimulus here
    wait;
  end process;
END;
{{ end }}
