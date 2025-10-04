"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion"
import { useEnvironmentMetadata } from "@/hooks/useRailway"
import { FileCode, GitBranch, Package, Clock, Settings, Box } from "lucide-react"
import { Skeleton } from "@/components/ui/skeleton"

interface EnvironmentMetadataProps {
  environmentId: string
  onSaveAsTemplate?: () => void
}

export function EnvironmentMetadata({ environmentId, onSaveAsTemplate }: EnvironmentMetadataProps) {
  const { data: metadata, isLoading, isError, error } = useEnvironmentMetadata(environmentId)

  if (isLoading) {
    return (
      <Card className="glass grain">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <Settings className="h-5 w-5" />
            Environment Configuration
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-3/4" />
          <Skeleton className="h-4 w-1/2" />
        </CardContent>
      </Card>
    )
  }

  if (isError || !metadata) {
    return (
      <Card className="glass grain border-amber-500/30">
        <CardHeader>
          <CardTitle className="text-lg flex items-center gap-2">
            <Settings className="h-5 w-5" />
            Environment Configuration
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            No configuration metadata available for this environment.
            {error && (
              <span className="block mt-2 text-xs text-amber-600">
                {error instanceof Error ? error.message : "Failed to load metadata"}
              </span>
            )}
          </p>
        </CardContent>
      </Card>
    )
  }

  const wizardInputs = metadata.wizardInputs || {}
  const provisionOutputs = metadata.provisionOutputs || {}

  return (
    <Card className="glass grain">
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg flex items-center gap-2">
            <Settings className="h-5 w-5" />
            Environment Configuration
          </CardTitle>
          {onSaveAsTemplate && (
            <Button variant="outline" size="sm" onClick={onSaveAsTemplate}>
              Save as Template
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent>
        <Accordion type="single" collapsible className="w-full">
          {/* Source Configuration */}
          {(wizardInputs.sourceType || wizardInputs.repositoryUrl || wizardInputs.dockerImage) && (
            <AccordionItem value="source">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  {wizardInputs.sourceType === "docker_image" ? (
                    <Box className="h-4 w-4 text-muted-foreground" />
                  ) : (
                    <GitBranch className="h-4 w-4 text-muted-foreground" />
                  )}
                  <span className="font-medium">Source</span>
                  <Badge variant="secondary" className="ml-2">
                    {wizardInputs.sourceType === "docker_image" ? "Docker Image" : "Repository"}
                  </Badge>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-3 pt-2">
                  {wizardInputs.sourceType === "repository" ? (
                    <>
                      {wizardInputs.repositoryUrl && (
                        <div>
                          <div className="text-xs text-muted-foreground mb-1">Repository URL</div>
                          <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                            {wizardInputs.repositoryUrl}
                          </div>
                        </div>
                      )}
                      {wizardInputs.branch && (
                        <div>
                          <div className="text-xs text-muted-foreground mb-1">Branch</div>
                          <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                            {wizardInputs.branch}
                          </div>
                        </div>
                      )}
                    </>
                  ) : (
                    <>
                      {wizardInputs.dockerImage && (
                        <div>
                          <div className="text-xs text-muted-foreground mb-1">Docker Image</div>
                          <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                            {wizardInputs.dockerImage}
                          </div>
                        </div>
                      )}
                      {wizardInputs.imageRegistry && (
                        <div>
                          <div className="text-xs text-muted-foreground mb-1">Registry</div>
                          <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                            {wizardInputs.imageRegistry}
                          </div>
                        </div>
                      )}
                      {wizardInputs.imageTag && (
                        <div>
                          <div className="text-xs text-muted-foreground mb-1">Tag</div>
                          <div className="text-sm font-mono bg-muted/30 px-3 py-2 rounded-md">
                            {wizardInputs.imageTag}
                          </div>
                        </div>
                      )}
                    </>
                  )}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}

          {/* Discovered Services */}
          {wizardInputs.discoveredServices && wizardInputs.discoveredServices.length > 0 && (
            <AccordionItem value="services">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Package className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Discovered Services</span>
                  <Badge variant="secondary" className="ml-2">
                    {wizardInputs.discoveredServices.length}
                  </Badge>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-2 pt-2">
                  {wizardInputs.discoveredServices.map((service, index) => (
                    <div
                      key={index}
                      className="bg-muted/30 px-3 py-2 rounded-md border border-border/50"
                    >
                      <div className="flex items-center justify-between">
                        <span className="text-sm font-medium">{service.name}</span>
                        {service.dockerfilePath && (
                          <FileCode className="h-3 w-3 text-muted-foreground" />
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground mt-1">
                        {service.path}
                      </div>
                      {service.dockerfilePath && (
                        <div className="text-xs text-muted-foreground mt-1 font-mono">
                          {service.dockerfilePath}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}

          {/* Environment Settings */}
          {(wizardInputs.environmentType || wizardInputs.ttl) && (
            <AccordionItem value="settings">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Settings className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Environment Settings</span>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-3 pt-2">
                  {wizardInputs.environmentType && (
                    <div>
                      <div className="text-xs text-muted-foreground mb-1">Environment Type</div>
                      <Badge variant="outline" className="capitalize">
                        {wizardInputs.environmentType}
                      </Badge>
                    </div>
                  )}
                  {wizardInputs.ttl && (
                    <div>
                      <div className="text-xs text-muted-foreground mb-1 flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        Time to Live (TTL)
                      </div>
                      <div className="text-sm bg-muted/30 px-3 py-2 rounded-md">
                        {wizardInputs.ttl} hours
                      </div>
                    </div>
                  )}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}

          {/* Provision Outputs */}
          {(provisionOutputs.railwayProjectId || provisionOutputs.railwayEnvironmentId) && (
            <AccordionItem value="outputs">
              <AccordionTrigger className="hover:no-underline">
                <div className="flex items-center gap-2">
                  <Box className="h-4 w-4 text-muted-foreground" />
                  <span className="font-medium">Provision Details</span>
                </div>
              </AccordionTrigger>
              <AccordionContent>
                <div className="space-y-3 pt-2">
                  {provisionOutputs.railwayProjectId && (
                    <div>
                      <div className="text-xs text-muted-foreground mb-1">Railway Project ID</div>
                      <div className="text-xs font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                        {provisionOutputs.railwayProjectId}
                      </div>
                    </div>
                  )}
                  {provisionOutputs.railwayEnvironmentId && (
                    <div>
                      <div className="text-xs text-muted-foreground mb-1">Railway Environment ID</div>
                      <div className="text-xs font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                        {provisionOutputs.railwayEnvironmentId}
                      </div>
                    </div>
                  )}
                  {provisionOutputs.serviceIds && Array.isArray(provisionOutputs.serviceIds) && (
                    <div>
                      <div className="text-xs text-muted-foreground mb-1">Service IDs</div>
                      <div className="text-xs font-mono bg-muted/30 px-3 py-2 rounded-md break-all">
                        {provisionOutputs.serviceIds.join(", ")}
                      </div>
                    </div>
                  )}
                </div>
              </AccordionContent>
            </AccordionItem>
          )}
        </Accordion>

        {/* Template Badge */}
        {metadata.isTemplate && (
          <div className="mt-4 pt-4 border-t border-border/50">
            <Badge variant="secondary" className="bg-blue-500/20 text-blue-700 border-blue-500/30">
              ⭐ Template
              {metadata.templateName && `: ${metadata.templateName}`}
            </Badge>
          </div>
        )}

        {/* Timestamps */}
        <div className="mt-4 pt-4 border-t border-border/50 text-xs text-muted-foreground">
          Created {new Date(metadata.createdAt).toLocaleString()}
        </div>
      </CardContent>
    </Card>
  )
}

