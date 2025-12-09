package com.example.curs3projectback.mapper;

import com.example.curs3projectback.dto.dictionary.AddressResponse;
import com.example.curs3projectback.dto.dictionary.OrganizationResponse;
import com.example.curs3projectback.model.Address;
import com.example.curs3projectback.model.Organization;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import java.util.List;

@Mapper(componentModel = "spring")
public interface DictionaryMapper {

    @Mapping(target = "consumer", expression = "java(address.getConsumer() != null ? address.getConsumer().getName() : null)")
    AddressResponse toAddress(Address address);

    List<AddressResponse> toAddressList(List<Address> addresses);

    OrganizationResponse toOrganization(Organization organization);

    List<OrganizationResponse> toOrganizationList(List<Organization> organizations);
}

