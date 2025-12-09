package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.ReplacementAct;
import com.example.curs3projectback.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;

public interface ReplacementActRepository extends JpaRepository<ReplacementAct, Long> {
    List<ReplacementAct> findAllByTask_Assignee(User assignee);
}

