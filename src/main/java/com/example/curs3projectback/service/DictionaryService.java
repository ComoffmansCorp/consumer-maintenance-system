package com.example.curs3projectback.service;

import com.example.curs3projectback.dto.dictionary.AddressResponse;
import com.example.curs3projectback.dto.dictionary.OrganizationResponse;
import com.example.curs3projectback.mapper.DictionaryMapper;
import com.example.curs3projectback.repository.AddressRepository;
import com.example.curs3projectback.repository.OrganizationRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class DictionaryService {

    private final AddressRepository addressRepository;
    private final OrganizationRepository organizationRepository;
    private final DictionaryMapper dictionaryMapper;

    public List<AddressResponse> findAddresses(String query) {
        return dictionaryMapper.toAddressList(addressRepository.findAllByStreetContainingIgnoreCaseOrderByStreet(query));
    }

    public List<OrganizationResponse> findOrganizations(String query) {
        return dictionaryMapper.toOrganizationList(organizationRepository.findAllByNameContainingIgnoreCaseOrderByName(query));
    }
}

