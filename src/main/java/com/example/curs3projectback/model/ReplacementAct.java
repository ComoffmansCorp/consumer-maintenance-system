package com.example.curs3projectback.model;

import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDate;
import java.util.ArrayList;
import java.util.List;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "replacement_acts")
public class ReplacementAct {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @OneToOne
    @JoinColumn(name = "task_id")
    private Task task;

    @ManyToOne(optional = false)
    private Address address;

    @Column(nullable = false)
    private String accountNumber;

    private LocalDate installationDate;

    @Embedded
    @AttributeOverrides({
            @AttributeOverride(name = "brand", column = @Column(name = "old_brand")),
            @AttributeOverride(name = "serialNumber", column = @Column(name = "old_serial_number")),
            @AttributeOverride(name = "readings", column = @Column(name = "old_readings"))
    })
    private MeterInfo oldMeter;

    @Embedded
    @AttributeOverrides({
            @AttributeOverride(name = "brand", column = @Column(name = "new_brand")),
            @AttributeOverride(name = "serialNumber", column = @Column(name = "new_serial_number")),
            @AttributeOverride(name = "readings", column = @Column(name = "new_readings"))
    })
    private MeterInfo newMeter;

    @OneToMany(mappedBy = "replacementAct", cascade = CascadeType.ALL, orphanRemoval = true)
    @Builder.Default
    private List<Photo> photos = new ArrayList<>();
}

