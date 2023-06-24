import { useState } from 'react'
import {
    makeStyles,
    Body1,
    Caption1,
    Button,
    shorthands,
    Card,
    CardFooter,
    CardHeader,
    CardPreview,
  } from "@fluentui/react-components"
import { CheckmarkRegular, DismissRegular } from "@fluentui/react-icons"
import Map from 'react-map-gl'
import './Client.css'
import api from './api'

function browserLogo(client) {
    const brand = client.info.browser_name.toLowerCase()
    return `https://raw.githubusercontent.com/alrra/browser-logos/main/src/${brand}/${brand}_128x128.png`
}

function osLogo(client) {
    const brand = client.info.os_name.toLowerCase()
    return `https://raw.githubusercontent.com/alrra/browser-logos/main/src/${brand}/${brand}_128x128.png`
}

const useStyles = makeStyles({
    card: {
      ...shorthands.margin("auto"),
      width: "720px",
      maxWidth: "100%",
    },
  })

function Client({client}) {
  console.log(client)

  const styles = useStyles();

  return (
    <Card className={styles.card}>
      <CardHeader
        image={
          <img src={browserLogo(client)} />
        }
        header={
          <Body1>
            <b>{client.info.ip}</b> tried to access <b>{client.resource}</b>
          </Body1>
        }
        description={<Caption1>Client ID: {client.id} Status: {client.status}</Caption1>}
      />
      <CardPreview>
        <div>
          <Body1>{client.info.city}, {client.info.region}, {client.info.country}</Body1>
        </div>
        <div>
          <Body1>{client.info.asn}</Body1>
        </div>
        <div>
          <div>
            <img src={browserLogo(client)} alt="" />
          </div>
          <div>
            <Body1>{client.info.browser_name} {client.info.browser_version}</Body1>
          </div>
          <div>
            <Body1>{client.info.os_name} {client.info.os_version}</Body1>
          </div>
          <div>
            <Map
              mapboxAccessToken='pk.eyJ1IjoiaWxvYWh6IiwiYSI6ImNqd2dlZDM3MDFlb3E0OG84OGptZmx4YTYifQ.or3xgAAaDIzk3TNpl0rfWQ'
              mapLib={import('mapbox-gl')}
              initialViewState={{
                longitude: client.info.long,
                latitude: client.info.lat,
                zoom: 6.5,
              }}
              style={{width: 600, height: 400}}
              mapStyle="mapbox://styles/mapbox/streets-v9"
            />
          </div>
        </div>
      </CardPreview>
      <CardFooter>
        <Button onClick={() => {api.allowClient(client.id)}} appearance='primary' icon={<CheckmarkRegular/>}>Allow</Button>
        <Button onClick={() => {api.blockClient(client.id)}} icon={<DismissRegular/>}>Block</Button>
      </CardFooter>
    </Card>
  )
}

export default Client
