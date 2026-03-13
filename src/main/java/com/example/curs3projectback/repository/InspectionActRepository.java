package com.example.curs3projectback.repository;

import com.example.curs3projectback.model.InspectionAct;
import com.example.curs3projectback.model.User;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.List;
import java.util.Optional;

public interface InspectionActRepository extends JpaRepository<InspectionAct, Long> {
    Optional<InspectionAct> findByIdAndTenant_Id(Long id, Long tenantId);
    List<InspectionAct> findAllByTask_Assignee(User assignee);
    List<InspectionAct> findAllByTask_AssigneeAndTenant_Id(User assignee, Long tenantId);
    long countByTenant_Id(Long tenantId);
}

