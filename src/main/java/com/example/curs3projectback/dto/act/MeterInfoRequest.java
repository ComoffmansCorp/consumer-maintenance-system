package com.example.curs3projectback.dto.act;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.PositiveOrZero;
import jakarta.validation.constraints.Size;
import lombok.Data;

@Data
public class MeterInfoRequest {
    @NotBlank(message = "brand is required")
    @Size(max = 255, message = "brand is too long")
    private String brand;
    @NotBlank(message = "serialNumber is required")
    @Size(max = 128, message = "serialNumber is too long")
    private String serialNumber;
    @NotNull(message = "readings is required")
    @PositiveOrZero(message = "readings must be >= 0")
    private Double readings;
}

