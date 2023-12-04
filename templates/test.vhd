LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;

ENTITY {{ .entityName }} IS
END {{ .entityName }};

ARCHITECTURE behavior OF {{ .entityName }} IS
  -- Component Declaration for the Unit Under Test (UUT)
  COMPONENT {{ .componentName }} IS
  PORT(
{{ .ports }}
  );
  END COMPONENT;
  --Inputs
{{ .inputs }}
  
  --Outputs
{{ .outputs }}

  -- Clock period definitions
  constant CLK50M_period : time := 20 ns;
BEGIN
  -- Instantiate the Unit Under Test (UUT)
  uut: {{ .componentName }} PORT MAP (
{{ .portMap }}
  );

  -- Clock process definitions
  CLK50M_process :process
  begin
		CLK50M <= '0';
		wait for CLK50M_period/2;
		CLK50M <= '1';
		wait for CLK50M_period/2;
  end process;

  -- Stimulus process
  stim_proc: process
  begin
    -- hold reset state for 100 ns.
	  RESET <= '1';
    wait for 100 ns;
    RESET <= '0';
    -- insert stimulus here
    wait;
  end process;
END;