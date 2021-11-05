#
# Copyright 2021. Clumio, Inc.
#
import collections.abc
import logging
import json
import os
from typing import Any

import boto3
from botocore import exceptions


"""Lambda event handler for acceptance test of clumio_callback_resource."""


class Error(Exception):
    """Base error class."""


class PropertyValueError(Error):
    """A required property is either undefined or empty."""


def clumio_event_handler(event: dict[str, Any], _: Any) -> None:
    """Event handler to process the clumio event."""
    logger = logging.Logger(name="process_event_logger")
    try:
        record = event['Records'][0]
        resource_properties = json.loads(record['Sns']['Message'])['ResourceProperties']
    except (KeyError, IndexError, TypeError, json.decoder.JSONDecodeError):
        raise Error(f'Unable to parse event for ResourceProperties')
    if not isinstance(resource_properties, dict):
        raise Error('ResourceProperties must be of type dictionary.')

    bucket_name = get_required_property(os.environ, 'BUCKET_NAME')
    account_id = get_required_property(resource_properties, 'AccountId')
    token = get_required_property(resource_properties, 'Token')
    region = get_required_property(resource_properties, 'Region')
    event_publish_time = get_required_property(resource_properties, 'EventPublishTime')
    canonical_user = get_required_property(resource_properties, 'CanonicalUser')
    client = boto3.client('s3')
    success_msg = {'Status': 'SUCCESS'}
    object_key = (
        f'acmtfstatus/{account_id}/{region}/{token}/clumio-status-{event_publish_time}.json')
    try:
        client.put_object(
            Bucket=bucket_name, Key=object_key, Body=bytes(json.dumps(success_msg), 'utf-8'))
        object_acl = client.get_object_acl(Bucket=bucket_name, Key=object_key)
        grants = object_acl.get('Grants', [])
        grants.append({
            'Grantee': {
                'ID': canonical_user,
                'Type': 'CanonicalUser'
            },
            'Permission': 'READ',
        })
        access_control_policy = {'Grants': grants,
                                 'Owner': object_acl.get('Owner')}
        client.put_object_acl(
            Bucket=bucket_name, Key=object_key, AccessControlPolicy=access_control_policy)
    except exceptions.ClientError as ex:
        logger.log(logging.ERROR, "Error in writing status message: %v", ex)
        raise Error("Error in writing status message.") from ex


def get_required_property(resource: collections.abc.Mapping, property_name: str) -> Any:
    """Raise error if required property is not present in the resource."""
    if not (data := resource.get(property_name)):
        raise PropertyValueError(f'Resource property "{property_name}" must be defined.')
    return data
