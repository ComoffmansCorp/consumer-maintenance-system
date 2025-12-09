package com.example.curs3projectback.model;

import com.example.curs3projectback.model.enums.InspectionType;
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
@Table(name = "inspection_acts")
public class InspectionAct {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @OneToOne
    @JoinColumn(name = "task_id")
    private Task task;

    @ManyToOne(optional = false)
    private Address address;

    private LocalDate inspectionDate;

    @ManyToOne
    private Organization consumer;

    @Enumerated(EnumType.STRING)
    private InspectionType inspectionType;

    @Column(length = 2000)
    private String notes;

    @OneToMany(mappedBy = "inspectionAct", cascade = CascadeType.ALL, orphanRemoval = true)
    @Builder.Default
    private List<Meter> meters = new ArrayList<>();

    @OneToMany(mappedBy = "inspectionAct", cascade = CascadeType.ALL, orphanRemoval = true)
    @Builder.Default
    private List<Photo> photos = new ArrayList<>();
}

