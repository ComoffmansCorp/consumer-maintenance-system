package com.example.curs3projectback.model;

import com.example.curs3projectback.model.enums.TaskStatus;
import com.example.curs3projectback.model.enums.TaskType;
import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDate;

@Getter
@Setter
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "tasks")
public class Task {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TaskType type;

    @ManyToOne(optional = false)
    private Address address;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private TaskStatus status;

    private LocalDate dueDate;

    @ManyToOne
    private User assignee;

    @OneToOne(mappedBy = "task")
    private InspectionAct inspectionAct;

    @OneToOne(mappedBy = "task")
    private ReplacementAct replacementAct;
}

