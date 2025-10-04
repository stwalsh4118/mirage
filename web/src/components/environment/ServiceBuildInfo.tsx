"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import { FileCode, GitBranch, Box, Settings, Boxes, Network } from "lucide-react"
import type { ServiceBuildConfig } from "@/lib/api/railway"

interface ServiceBuildInfoProps {
  service: ServiceBuildConfig
  compact?: boolean
}

export function ServiceBuildInfo({ service, compact = false }: ServiceBuildInfoProps) {
  const isRepositoryBased = service.deploymentType === "repository"
  const hasBuildArgs = service.buildArgs && Object.keys(service.buildArgs).length > 0
  const hasPorts = service.exposedPorts && service.exposedPorts.length > 0

  if (compact) {
    return (
      <div className="space-y-2 text-sm">
        {/* Deployment Type */}
        <div className="flex items-center gap-2">
          {isRepositoryBased ? (
            <GitBranch className="h-4 w-4 text-muted-foreground" />
          ) : (
            <Box className="h-4 w-4 text-muted-foreground" />
          )}
          <Badge variant="outline" className="text-xs">
            {isRepositoryBased ? "Repository" : "Docker Image"}
          </Badge>
        </div>

        {/* Source/Image */}
        {isRepositoryBased ? (
          <div>
            {service.sourceRepo && (
              <div className="text-xs text-muted-foreground break-all">
                {service.sourceRepo}
                {service.sourceBranch && ` @ ${service.sourceBranch}`}
              </div>
            )}
            {service.dockerfilePath && (
              <div className="text-xs font-mono text-muted-foreground mt-1">
                <FileCode className="h-3 w-3 inline mr-1" />
                {service.dockerfilePath}
              </div>
            )}
          </div>
        ) : (
          <div>
            {service.dockerImage && (
              <div className="text-xs font-mono break-all text-muted-foreground">
                {service.imageRegistry && service.imageRegistry !== "docker.io" && `${service.imageRegistry}/`}
                {service.dockerImage}
                {service.imageTag && `:${service.imageTag}`}
              </div>
            )}
          </div>
        )}

        {/* Ports */}
        {hasPorts && (
          <div className="flex items-center gap-1">
            <Network className="h-3 w-3 text-muted-foreground" />
            <span className="text-xs text-muted-foreground">
              Ports: {service.exposedPorts?.join(", ")}
            </span>
          </div>
        )}
      </div>
    )
  }

  return (
    <Card className="glass grain">
      <CardHeader className="pb-3">
        <CardTitle className="text-base flex items-center gap-2">
          <Boxes className="h-4 w-4" />
          {service.name}
          <Badge variant="secondary" className="ml-auto">
            {service.status}
          </Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <Accordion type="single" collapsible className="w-full">
          {/* Deployment Configuration */}
          <AccordionItem value="deployment">
            <AccordionTrigger className="hover:no-underline">
              <div className="flex items-center gap-2">
                {isRepositoryBased ? (
                  <GitBranch className="h-4 w-4 text-muted-foreground" />
                ) : (
                  <Box className="h-4 w-4 text-muted-foreground" />
                )}
                <span className="font-medium">Deployment</span>
                <Badge variant="outline" className="ml-2">
                  {isRepositoryBased ? "Repository" : "Docker Image"}
                </Badge>
              </div>
            </AccordionTrigger>
            <AccordionContent>
              <div className="space-y-3 pt-2">
                {isRepositoryBased ? (
                  <>
                    {service.sourceRepo && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Repository</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                          {service.sourceRepo}
                        </div>
                      </div>
                    )}
                    {service.sourceBranch && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Branch</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                          {service.sourceBranch}
                        </div>
                      </div>
                    )}
                    {service.dockerfilePath && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1 flex items-center gap-1">
                          <FileCode className="h-3 w-3" />
                          Dockerfile Path
                        </div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                          {service.dockerfilePath}
                        </div>
                      </div>
                    )}
                    {service.buildContext && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Build Context</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                          {service.buildContext}
                        </div>
                      </div>
                    )}
                  </>
                ) : (
                  <>
                    {service.dockerImage && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Docker Image</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                          {service.dockerImage}
                        </div>
                      </div>
                    )}
                    {service.imageRegistry && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Registry</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                          {service.imageRegistry}
                        </div>
                      </div>
                    )}
                    {service.imageTag && (
                      <div>
                        <div className="text-xs text-muted-foreground mb-1">Tag</div>
                        <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                          {service.imageTag}
                        </div>
                      </div>
                    )}
                  </>
                )}
              </div>
            </AccordionContent>
          </AccordionItem>

          {/* Build Configuration */}
          {hasBuildArgs && (
            <AccordionItem value="build">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Settings className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Build Arguments</span>
                  <Badge variant="secondary" className="ml-2">
                    {Object.keys(service.buildArgs!).length}
                  </Badge>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-2 pt-2">
                  {Object.entries(service.buildArgs!).map(([key, value]) => (
                    <div key={key} className="bg-muted/30 px-3 py-2 rounded-md border border-border/50">
                      <div className="text-xs text-muted-foreground mb-1">{key}</div>
                      <div className="text-sm font-mono break-all">{value}</div>
                    </div>
                  ))}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}

          {/* Network Configuration */}
          {hasPorts && (
            <AccordionItem value="network">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Network className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Exposed Ports</span>
                  <Badge variant="secondary" className="ml-2">
                    {service.exposedPorts!.length}
                  </Badge>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="pt-2 flex flex-wrap gap-2">
                  {service.exposedPorts!.map((port) => (
                    <Badge key={port} variant="outline" className="font-mono">
                      {port}
                    </Badge>
                  ))}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}

          {/* Railway Service ID */}
          {service.railwayServiceId && (
            <AccordionItem value="railway">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Box className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Railway Details</span>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="pt-2">
                  <div className="text-xs text-muted-foreground mb-1">Service ID</div>
                  <div className="text-xs font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                    {service.railwayServiceId}
                  </div>
                </div>
              </AccordionContent>
            </AccordionItem>
          )}
        </Accordion>

        {/* Timestamp */}
        <div className="mt-4 pt-4 border-t border-border/50 text-xs text-muted-foreground">
          Created {new Date(service.createdAt).toLocaleString()}
        </div>
      </CardContent>
    </Card>
  )
}

