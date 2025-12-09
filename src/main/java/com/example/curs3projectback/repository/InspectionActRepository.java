package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.InspectionAct;
import com.example.curs3projectback.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface InspectionActRepository extends JpaRepository<InspectionAct, Long> {
    List<InspectionAct> findAllByTask_Assignee(User assignee);
}

