package io.github.modernstack.forge;

import org.springframework.boot.SpringApplication;

public class TestDaggerTestcontainersForgeApplication {

	public static void main(String[] args) {
		SpringApplication.from(DaggerTestcontainersForgeApplication::main).with(TestcontainersConfiguration.class).run(args);
	}

}
