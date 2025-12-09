package com.example.curs3projectback.model;

import com.example.curs3projectback.model.enums.MeterType;
import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDate;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "meters")
public class Meter {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private MeterType type;

    @Column(nullable = false)
    private String serialNumber;

    private Integer manufactureYear;

    private LocalDate verificationDate;

    private String sealState;

    private Integer transformationRatio;

    @ManyToOne
    @JoinColumn(name = "inspection_act_id")
    private InspectionAct inspectionAct;
}

