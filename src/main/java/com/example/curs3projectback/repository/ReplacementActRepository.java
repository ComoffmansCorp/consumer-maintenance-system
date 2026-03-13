package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.ReplacementAct;
import com.example.curs3projectback.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface ReplacementActRepository extends JpaRepository<ReplacementAct, Long> {
    Optional<ReplacementAct> findByIdAndTenant_Id(Long id, Long tenantId);
    List<ReplacementAct> findAllByTask_Assignee(User assignee);
    List<ReplacementAct> findAllByTask_AssigneeAndTenant_Id(User assignee, Long tenantId);
    long countByTenant_Id(Long tenantId);
}

