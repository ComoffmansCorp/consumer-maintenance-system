package com.example.curs3projectback.model;

import jakarta.persistence.Embeddable;
import lombok.*;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Embeddable
public class MeterInfo {
    private String brand;
    private String serialNumber;
    private Double readings;
}

