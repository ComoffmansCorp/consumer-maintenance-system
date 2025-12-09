package com.example.curs3projectback.model;

import jakarta.persistence.*;
import lombok.*;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "photos")
public class Photo {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false)
    private String filename;

    private String note;

    @ManyToOne
    @JoinColumn(name = "inspection_act_id")
    private InspectionAct inspectionAct;

    @ManyToOne
    @JoinColumn(name = "replacement_act_id")
    private ReplacementAct replacementAct;
}

